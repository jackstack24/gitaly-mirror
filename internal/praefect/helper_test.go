package praefect

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/client"
	gconfig "gitlab.com/gitlab-org/gitaly/internal/config"
	internalauth "gitlab.com/gitlab-org/gitaly/internal/config/auth"
	"gitlab.com/gitlab-org/gitaly/internal/log"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/grpc-proxy/proxy"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/mock"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/models"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/nodes"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/protoregistry"
	"gitlab.com/gitlab-org/gitaly/internal/server/auth"
	"gitlab.com/gitlab-org/gitaly/internal/service/internalgitaly"
	"gitlab.com/gitlab-org/gitaly/internal/service/repository"
	gitalyserver "gitlab.com/gitlab-org/gitaly/internal/service/server"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper/promtest"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	correlation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func waitUntil(t *testing.T, ch <-chan struct{}, timeout time.Duration) {
	select {
	case <-ch:
		break
	case <-time.After(timeout):
		t.Errorf("timed out waiting for channel after %s", timeout)
	}
}

// generates a praefect configuration with the specified number of backend
// nodes
func testConfig(backends int) config.Config {
	var nodes []*models.Node

	for i := 0; i < backends; i++ {
		n := &models.Node{
			Storage: fmt.Sprintf("praefect-internal-%d", i),
			Token:   fmt.Sprintf("%d", i),
		}

		if i == 0 {
			n.DefaultPrimary = true
		}

		nodes = append(nodes, n)
	}
	cfg := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			&config.VirtualStorage{
				Name:  "praefect",
				Nodes: nodes,
			},
		},
	}

	return cfg
}

// setupServer wires all praefect dependencies together via dependency
// injection
func setupServer(t testing.TB, conf config.Config, nodeMgr nodes.Manager, ds datastore.Datastore, l *logrus.Entry, r *protoregistry.Registry) *Server {
	coordinator := NewCoordinator(l, ds, nodeMgr, conf, r)

	var defaultNode *models.Node
	for _, n := range conf.VirtualStorages[0].Nodes {
		if n.DefaultPrimary {
			defaultNode = n
		}
	}
	require.NotNil(t, defaultNode)

	server := NewServer(coordinator.StreamDirector, l, r, conf)

	return server
}

// runPraefectServer runs a praefect server with the provided mock servers.
// Each mock server is keyed by the corresponding index of the node in the
// config.Nodes. There must be a 1-to-1 mapping between backend server and
// configured storage node.
// requires there to be only 1 virtual storage
func runPraefectServerWithMock(t *testing.T, conf config.Config, ds datastore.Datastore, backends map[string]mock.SimpleServiceServer) (*grpc.ClientConn, *Server, testhelper.Cleanup) {
	require.Len(t, conf.VirtualStorages, 1)
	require.Equal(t, len(backends), len(conf.VirtualStorages[0].Nodes),
		"mock server count doesn't match config nodes")

	var cleanups []testhelper.Cleanup

	for i, node := range conf.VirtualStorages[0].Nodes {
		backend, ok := backends[node.Storage]
		require.True(t, ok, "missing backend server for node %s", node.Storage)

		backendAddr, cleanup := newMockDownstream(t, node.Token, backend)
		cleanups = append(cleanups, cleanup)

		node.Address = backendAddr
		conf.VirtualStorages[0].Nodes[i] = node
	}

	nodeMgr, err := nodes.NewManager(testhelper.DiscardTestEntry(t), conf, nil, datastore.Datastore{}, promtest.NewMockHistogramVec())
	require.NoError(t, err)
	nodeMgr.Start(1*time.Millisecond, 5*time.Millisecond)

	r := protoregistry.New()
	require.NoError(t, r.RegisterFiles(mustLoadProtoReg(t)))

	prf := setupServer(t, conf, nodeMgr, ds, log.Default(), r)

	listener, port := listenAvailPort(t)
	t.Logf("praefect listening on port %d", port)

	errQ := make(chan error)

	prf.RegisterServices(nodeMgr, conf, ds)
	go func() {
		errQ <- prf.Serve(listener, false)
	}()

	// dial client to praefect
	cc := dialLocalPort(t, port, false)

	cleanup := func() {
		for _, cu := range cleanups {
			cu()
		}
		require.NoError(t, cc.Close())
		require.NoError(t, listener.Close())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		require.NoError(t, prf.Shutdown(ctx))
	}

	return cc, prf, cleanup
}

func noopBackoffFunc() (backoff, backoffReset) {
	return func() time.Duration {
		return 0
	}, func() {}
}

func runPraefectServerWithGitaly(t *testing.T, conf config.Config) (*grpc.ClientConn, *Server, testhelper.Cleanup) {
	ds := datastore.Datastore{
		ReplicasDatastore:     datastore.NewInMemory(conf),
		ReplicationEventQueue: datastore.NewMemoryReplicationEventQueue(conf),
	}

	return runPraefectServerWithGitalyWithDatastore(t, conf, ds)
}

// runPraefectServerWithGitaly runs a praefect server with actual Gitaly nodes
// requires exactly 1 virtual storage
func runPraefectServerWithGitalyWithDatastore(t *testing.T, conf config.Config, ds datastore.Datastore) (*grpc.ClientConn, *Server, testhelper.Cleanup) {
	require.Len(t, conf.VirtualStorages, 1)
	var cleanups []testhelper.Cleanup

	var storages []gconfig.Storage
	for _, node := range conf.VirtualStorages[0].Nodes {
		storages = append(storages, gconfig.Storage{
			Name: node.Storage,
			Path: testhelper.GitlabTestStoragePath(),
		})
	}

	_, backendAddr, cleanupGitaly := runInternalGitalyServer(t, storages, conf.VirtualStorages[0].Nodes[0].Token)
	cleanups = append(cleanups, cleanupGitaly)

	for i, node := range conf.VirtualStorages[0].Nodes {
		node.Address = backendAddr
		conf.VirtualStorages[0].Nodes[i] = node
	}

	logEntry := log.Default()

	nodeMgr, err := nodes.NewManager(testhelper.DiscardTestEntry(t), conf, nil, ds, promtest.NewMockHistogramVec())
	require.NoError(t, err)
	nodeMgr.Start(1*time.Millisecond, 5*time.Millisecond)

	registry := protoregistry.New()
	require.NoError(t, registry.RegisterFiles(protoregistry.GitalyProtoFileDescriptors...))
	coordinator := NewCoordinator(logEntry, ds, nodeMgr, conf, registry)

	replmgr := NewReplMgr(
		conf.VirtualStorages[0].Name,
		logEntry,
		ds,
		nodeMgr,
		WithQueueMetric(&promtest.MockGauge{}),
	)
	prf := NewServer(
		coordinator.StreamDirector,
		logEntry,
		registry,
		conf,
	)

	listener, port := listenAvailPort(t)
	t.Logf("proxy listening on port %d", port)

	errQ := make(chan error)
	ctx, cancel := testhelper.Context()

	prf.RegisterServices(nodeMgr, conf, ds)
	go func() { errQ <- prf.Serve(listener, false) }()
	go func() { replmgr.ProcessBacklog(ctx, noopBackoffFunc) }()

	// dial client to praefect
	cc := dialLocalPort(t, port, false)

	cleanup := func() {
		for _, cu := range cleanups {
			cu()
		}

		ctx, timed := context.WithTimeout(ctx, time.Second)
		defer timed()
		require.NoError(t, prf.Shutdown(ctx))

		cancel()
		require.Error(t, context.Canceled, <-errQ)
	}

	return cc, prf, cleanup
}

func runInternalGitalyServer(t *testing.T, storages []gconfig.Storage, token string) (*grpc.Server, string, func()) {
	streamInt := []grpc.StreamServerInterceptor{auth.StreamServerInterceptor(internalauth.Config{Token: token})}
	unaryInt := []grpc.UnaryServerInterceptor{auth.UnaryServerInterceptor(internalauth.Config{Token: token})}

	server := testhelper.NewTestGrpcServer(t, streamInt, unaryInt)
	serverSocketPath := testhelper.GetTemporaryGitalySocketFileName()

	listener, err := net.Listen("unix", serverSocketPath)
	require.NoError(t, err)

	internalSocket := gconfig.GitalyInternalSocketPath()
	internalListener, err := net.Listen("unix", internalSocket)
	require.NoError(t, err)

	gitalypb.RegisterServerServiceServer(server, gitalyserver.NewServer(storages))
	gitalypb.RegisterRepositoryServiceServer(server, repository.NewServer(RubyServer, internalSocket))
	gitalypb.RegisterInternalGitalyServer(server, internalgitaly.NewServer(gconfig.Config.Storages))
	healthpb.RegisterHealthServer(server, health.NewServer())

	errQ := make(chan error)

	go func() {
		errQ <- server.Serve(listener)
	}()
	go func() {
		errQ <- server.Serve(internalListener)
	}()

	cleanup := func() {
		server.Stop()
		require.NoError(t, <-errQ)
	}

	return server, "unix://" + serverSocketPath, cleanup
}

func mustLoadProtoReg(t testing.TB) *descriptor.FileDescriptorProto {
	gz, _ := (*mock.SimpleRequest)(nil).Descriptor()
	fd, err := protoregistry.ExtractFileDescriptor(gz)
	require.NoError(t, err)
	return fd
}

func listenAvailPort(tb testing.TB) (net.Listener, int) {
	listener, err := net.Listen("tcp", ":0")
	require.NoError(tb, err)

	return listener, listener.Addr().(*net.TCPAddr).Port
}

func dialLocalPort(tb testing.TB, port int, backend bool) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(correlation.UnaryClientCorrelationInterceptor()),
		grpc.WithStreamInterceptor(correlation.StreamClientCorrelationInterceptor()),
	}
	if backend {
		opts = append(
			opts,
			grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.NewCodec())),
		)
	}

	cc, err := client.Dial(
		fmt.Sprintf("tcp://localhost:%d", port),
		opts,
	)
	require.NoError(tb, err)

	return cc
}

func newMockDownstream(tb testing.TB, token string, m mock.SimpleServiceServer) (string, func()) {
	srv := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryServerInterceptor(internalauth.Config{Token: token})))
	mock.RegisterSimpleServiceServer(srv, m)

	// client to backend service
	lis, port := listenAvailPort(tb)

	errQ := make(chan error)

	go func() {
		errQ <- srv.Serve(lis)
	}()

	cleanup := func() {
		srv.GracefulStop()
		lis.Close()

		// If the server is shutdown before Serve() is called on it
		// the Serve() calls will return the ErrServerStopped
		if err := <-errQ; err != nil && err != grpc.ErrServerStopped {
			require.NoError(tb, err)
		}
	}

	return fmt.Sprintf("tcp://localhost:%d", port), cleanup
}
