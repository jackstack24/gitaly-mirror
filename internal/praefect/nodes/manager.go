package nodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	gitalyauth "gitlab.com/gitlab-org/gitaly/auth"
	"gitlab.com/gitlab-org/gitaly/client"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/grpc-proxy/proxy"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/metrics"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/models"
	prommetrics "gitlab.com/gitlab-org/gitaly/internal/prometheus/metrics"
	correlation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	grpctracing "gitlab.com/gitlab-org/labkit/tracing/grpc"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Shard is a primary with a set of secondaries
type Shard interface {
	GetPrimary() (Node, error)
	GetSecondaries() ([]Node, error)
}

// Manager is responsible for returning shards for virtual storages
type Manager interface {
	GetShard(virtualStorageName string) (Shard, error)
	// GetSyncedNode returns a random storage node based on the state of the replication.
	// It returns primary in case there are no up to date secondaries or error occurs.
	GetSyncedNode(ctx context.Context, virtualStorageName, repoPath string) (Node, error)
}

// Node represents some metadata of a node as well as a connection
type Node interface {
	GetStorage() string
	GetAddress() string
	GetToken() string
	GetConnection() *grpc.ClientConn
	check(context.Context) (bool, error)
}

// Mgr is a concrete type that adheres to the Manager interface
type Mgr struct {
	failoverEnabled bool
	// log must be used only for non-request specific needs like bootstrapping, etc.
	// for request related logging `ctxlogrus.Extract(ctx)` must be used.
	log *logrus.Entry
	// strategies is a map of strategies keyed on virtual storage name
	strategies map[string]leaderElectionStrategy
	db         *sql.DB
	ds         datastore.Datastore
}

// leaderElectionStrategy defines the interface by which primary and
// secondaries are managed.
type leaderElectionStrategy interface {
	start(bootstrapInterval, monitorInterval time.Duration)
	checkNodes(context.Context) error

	Shard
}

// ErrPrimaryNotHealthy indicates the primary of a shard is not in a healthy state and hence
// should not be used for a new request
var ErrPrimaryNotHealthy = errors.New("primary is not healthy")

// NewManager creates a new NodeMgr based on virtual storage configs
func NewManager(log *logrus.Entry, c config.Config, db *sql.DB, ds datastore.Datastore, latencyHistogram prommetrics.HistogramVec, dialOpts ...grpc.DialOption) (*Mgr, error) {
	strategies := make(map[string]leaderElectionStrategy, len(c.VirtualStorages))

	for _, virtualStorage := range c.VirtualStorages {
		ns := make([]*nodeStatus, 1, len(virtualStorage.Nodes))
		for _, node := range virtualStorage.Nodes {
			conn, err := client.Dial(node.Address,
				append(
					[]grpc.DialOption{
						grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.NewCodec())),
						grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2(node.Token)),
						grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
							grpc_prometheus.StreamClientInterceptor,
							grpctracing.StreamClientTracingInterceptor(),
							correlation.StreamClientCorrelationInterceptor(),
						)),
						grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
							grpc_prometheus.UnaryClientInterceptor,
							grpctracing.UnaryClientTracingInterceptor(),
							correlation.UnaryClientCorrelationInterceptor(),
						)),
					}, dialOpts...),
			)
			if err != nil {
				return nil, err
			}
			cs := newConnectionStatus(*node, conn, log, latencyHistogram)
			if node.DefaultPrimary {
				ns[0] = cs
			} else {
				ns = append(ns, cs)
			}
		}

		if c.Failover.ElectionStrategy == "sql" {
			strategies[virtualStorage.Name] = newSQLElector(virtualStorage.Name, c, defaultFailoverTimeoutSeconds, defaultActivePraefectSeconds, db, log, ns)
		} else {
			strategies[virtualStorage.Name] = newLocalElector(virtualStorage.Name, c.Failover.Enabled, log, ns)
		}
	}

	return &Mgr{
		log:             log,
		db:              db,
		failoverEnabled: c.Failover.Enabled,
		strategies:      strategies,
		ds:              ds,
	}, nil
}

// Start will bootstrap the node manager by calling healthcheck on the nodes as well as kicking off
// the monitoring process. Start must be called before NodeMgr can be used.
func (n *Mgr) Start(bootstrapInterval, monitorInterval time.Duration) {
	if n.failoverEnabled {
		n.log.Info("Starting failover checks")

		for _, strategy := range n.strategies {
			strategy.start(bootstrapInterval, monitorInterval)
		}
	} else {
		n.log.Info("Failover checks are disabled")
	}
}

// checkShards performs health checks on all the available shards. The
// election strategy is responsible for determining the criteria for
// when to elect a new primary and when a node is down.
func (n *Mgr) checkShards() {
	for _, strategy := range n.strategies {
		ctx := context.Background()
		strategy.checkNodes(ctx)
	}
}

// ErrVirtualStorageNotExist indicates the node manager is not aware of the virtual storage for which a shard is being requested
var ErrVirtualStorageNotExist = errors.New("virtual storage does not exist")

// GetShard retrieves a shard for a virtual storage name
func (n *Mgr) GetShard(virtualStorageName string) (Shard, error) {
	shard, ok := n.strategies[virtualStorageName]
	if !ok {
		return nil, ErrVirtualStorageNotExist
	}

	if n.failoverEnabled {
		_, err := shard.GetPrimary()

		if err != nil {
			return nil, ErrPrimaryNotHealthy
		}
	}

	return shard, nil
}

func (n *Mgr) GetSyncedNode(ctx context.Context, virtualStorageName, repoPath string) (Node, error) {
	shard, err := n.GetShard(virtualStorageName)
	if err != nil {
		return nil, fmt.Errorf("get shard for %q: %w", virtualStorageName, err)
	}

	storages, err := n.ds.GetUpToDateStorages(ctx, virtualStorageName, repoPath)
	if err != nil {
		// this is recoverable error - proceed with primary node
		ctxlogrus.Extract(ctx).
			WithError(err).
			WithFields(logrus.Fields{"virtual_storage_name": virtualStorageName, "repo_path": repoPath}).
			Error("get up to date secondaries")
	}

	if len(storages) == 0 {
		primary, err := shard.GetPrimary()
		if err != nil {
			return nil, fmt.Errorf("no synced, get primary for %q: %w", virtualStorageName, err)
		}
		return primary, nil
	}

	secondaries, err := shard.GetSecondaries()
	if err != nil {
		return nil, fmt.Errorf("get secondaries for %q: %w", virtualStorageName, err)
	}

	secondary := storages[rand.Intn(len(storages))] // randomly pick up one of the synced storages
	for _, node := range secondaries {
		if node.GetStorage() == secondary {
			return node, nil
		}
	}

	primary, err := shard.GetPrimary() // there is no matched secondaries, maybe because of re-configuration
	if err != nil {
		return nil, fmt.Errorf("get primary for %q: %w", virtualStorageName, err)
	}
	return primary, nil
}

func newConnectionStatus(node models.Node, cc *grpc.ClientConn, l *logrus.Entry, latencyHist prommetrics.HistogramVec) *nodeStatus {
	return &nodeStatus{
		Node:        node,
		ClientConn:  cc,
		log:         l,
		latencyHist: latencyHist,
	}
}

type nodeStatus struct {
	models.Node
	*grpc.ClientConn
	log         *logrus.Entry
	latencyHist prommetrics.HistogramVec
}

// GetStorage gets the storage name of a node
func (n *nodeStatus) GetStorage() string {
	return n.Storage
}

// GetAddress gets the address of a node
func (n *nodeStatus) GetAddress() string {
	return n.Address
}

// GetToken gets the token of a node
func (n *nodeStatus) GetToken() string {
	return n.Token
}

// GetConnection gets the client connection of a node
func (n *nodeStatus) GetConnection() *grpc.ClientConn {
	return n.ClientConn
}

func (n *nodeStatus) check(ctx context.Context) (bool, error) {
	client := healthpb.NewHealthClient(n.ClientConn)
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	status := false

	start := time.Now()
	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
	n.latencyHist.WithLabelValues(n.Storage).Observe(time.Since(start).Seconds())

	if err == nil && resp.Status == healthpb.HealthCheckResponse_SERVING {
		status = true
	} else {
		n.log.WithError(err).WithFields(logrus.Fields{
			"storage": n.Storage,
			"address": n.Address,
		}).Warn("error when pinging healthcheck")
	}

	var gaugeValue float64
	if status {
		gaugeValue = 1
	}
	metrics.NodeLastHealthcheckGauge.WithLabelValues(n.GetStorage()).Set(gaugeValue)

	return status, err
}
