package nodes

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/models"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper/promtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestNodeStatus(t *testing.T) {
	socket := testhelper.GetTemporaryGitalySocketFileName()
	svr, healthSvr := testhelper.NewServerWithHealth(t, socket)
	defer svr.Stop()

	cc, err := grpc.Dial(
		"unix://"+socket,
		grpc.WithInsecure(),
	)

	require.NoError(t, err)

	mockHistogramVec := promtest.NewMockHistogramVec()

	storageName := "default"
	cs := newConnectionStatus(models.Node{Storage: storageName}, cc, testhelper.DiscardTestEntry(t), mockHistogramVec)

	var expectedLabels [][]string
	for i := 0; i < healthcheckThreshold; i++ {
		ctx := context.Background()
		status, err := cs.check(ctx)

		require.NoError(t, err)
		require.True(t, status)
		expectedLabels = append(expectedLabels, []string{storageName})
	}

	require.Equal(t, expectedLabels, mockHistogramVec.LabelsCalled())
	require.Len(t, mockHistogramVec.Observer().Observed(), healthcheckThreshold)

	healthSvr.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	ctx := context.Background()
	status, err := cs.check(ctx)
	require.NoError(t, err)
	require.False(t, status)
}

func TestPrimaryIsSecond(t *testing.T) {
	virtualStorages := []*config.VirtualStorage{
		{
			Name: "virtual-storage-0",
			Nodes: []*models.Node{
				{
					Storage:        "praefect-internal-0",
					Address:        "unix://socket0",
					DefaultPrimary: false,
				},
				{
					Storage:        "praefect-internal-1",
					Address:        "unix://socket1",
					DefaultPrimary: true,
				},
			},
		},
	}

	conf := config.Config{
		VirtualStorages: virtualStorages,
		Failover:        config.Failover{Enabled: false},
	}

	mockHistogram := promtest.NewMockHistogramVec()
	nm, err := NewManager(testhelper.DiscardTestEntry(t), conf, nil, datastore.Datastore{}, mockHistogram)
	require.NoError(t, err)

	shard, err := nm.GetShard("virtual-storage-0")
	require.NoError(t, err)

	primary, err := shard.GetPrimary()
	require.NoError(t, err)

	secondaries, err := shard.GetSecondaries()
	require.Len(t, secondaries, 1)
	require.NoError(t, err)

	require.Equal(t, virtualStorages[0].Nodes[1].Storage, primary.GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[1].Address, primary.GetAddress())

	require.Len(t, secondaries, 1)
	require.Equal(t, virtualStorages[0].Nodes[0].Storage, secondaries[0].GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[0].Address, secondaries[0].GetAddress())
}

func TestNodeManager(t *testing.T) {
	internalSocket0, internalSocket1 := testhelper.GetTemporaryGitalySocketFileName(), testhelper.GetTemporaryGitalySocketFileName()
	srv0, healthSrv0 := testhelper.NewServerWithHealth(t, internalSocket0)
	defer srv0.Stop()

	srv1, healthSrv1 := testhelper.NewServerWithHealth(t, internalSocket1)
	defer srv1.Stop()

	virtualStorages := []*config.VirtualStorage{
		{
			Name: "virtual-storage-0",
			Nodes: []*models.Node{
				{
					Storage:        "praefect-internal-0",
					Address:        "unix://" + internalSocket0,
					DefaultPrimary: true,
				},
				{
					Storage: "praefect-internal-1",
					Address: "unix://" + internalSocket1,
				},
			},
		},
	}

	confWithFailover := config.Config{
		VirtualStorages: virtualStorages,
		Failover:        config.Failover{Enabled: true},
	}
	confWithoutFailover := config.Config{
		VirtualStorages: virtualStorages,
		Failover:        config.Failover{Enabled: false},
	}

	ds := datastore.Datastore{}

	mockHistogram := promtest.NewMockHistogramVec()
	nm, err := NewManager(testhelper.DiscardTestEntry(t), confWithFailover, nil, ds, mockHistogram)
	require.NoError(t, err)

	nmWithoutFailover, err := NewManager(testhelper.DiscardTestEntry(t), confWithoutFailover, nil, ds, mockHistogram)
	require.NoError(t, err)

	nm.Start(1*time.Millisecond, 5*time.Second)
	nmWithoutFailover.Start(1*time.Millisecond, 5*time.Second)

	_, err = nm.GetShard("virtual-storage-0")
	require.NoError(t, err)

	shardWithoutFailover, err := nmWithoutFailover.GetShard("virtual-storage-0")
	require.NoError(t, err)
	primaryWithoutFailover, err := shardWithoutFailover.GetPrimary()
	require.NoError(t, err)
	secondariesWithoutFailover, err := shardWithoutFailover.GetSecondaries()
	require.NoError(t, err)

	shard, err := nm.GetShard("virtual-storage-0")
	require.NoError(t, err)
	primary, err := shard.GetPrimary()
	require.NoError(t, err)
	secondaries, err := shard.GetSecondaries()
	require.NoError(t, err)

	// shard without failover and shard with failover should be the same
	require.Equal(t, primaryWithoutFailover.GetStorage(), primary.GetStorage())
	require.Equal(t, primaryWithoutFailover.GetAddress(), primary.GetAddress())
	require.Len(t, secondaries, 1)
	require.Equal(t, secondariesWithoutFailover[0].GetStorage(), secondaries[0].GetStorage())
	require.Equal(t, secondariesWithoutFailover[0].GetAddress(), secondaries[0].GetAddress())

	require.Equal(t, virtualStorages[0].Nodes[0].Storage, primary.GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[0].Address, primary.GetAddress())
	require.Len(t, secondaries, 1)
	require.Equal(t, virtualStorages[0].Nodes[1].Storage, secondaries[0].GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[1].Address, secondaries[0].GetAddress())

	healthSrv0.SetServingStatus("", grpc_health_v1.HealthCheckResponse_UNKNOWN)
	nm.checkShards()

	labelsCalled := mockHistogram.LabelsCalled()
	for _, node := range virtualStorages[0].Nodes {
		require.Contains(t, labelsCalled, []string{node.Storage})
	}

	// since the primary is unhealthy, we expect checkShards to demote primary to secondary, and promote the healthy
	// secondary to primary

	shardWithoutFailover, err = nmWithoutFailover.GetShard("virtual-storage-0")
	require.NoError(t, err)
	primaryWithoutFailover, err = shardWithoutFailover.GetPrimary()
	require.NoError(t, err)
	secondariesWithoutFailover, err = shardWithoutFailover.GetSecondaries()
	require.NoError(t, err)

	shard, err = nm.GetShard("virtual-storage-0")
	require.NoError(t, err)
	primary, err = shard.GetPrimary()
	require.NoError(t, err)
	secondaries, err = shard.GetSecondaries()
	require.NoError(t, err)

	// shard without failover and shard with failover should not be the same
	require.NotEqual(t, primaryWithoutFailover.GetStorage(), primary.GetStorage())
	require.NotEqual(t, primaryWithoutFailover.GetAddress(), primary.GetAddress())
	require.NotEqual(t, secondariesWithoutFailover[0].GetStorage(), secondaries[0].GetStorage())
	require.NotEqual(t, secondariesWithoutFailover[0].GetAddress(), secondaries[0].GetAddress())

	// shard without failover should still match the config
	require.Equal(t, virtualStorages[0].Nodes[0].Storage, primaryWithoutFailover.GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[0].Address, primaryWithoutFailover.GetAddress())
	require.Len(t, secondaries, 1)
	require.Equal(t, virtualStorages[0].Nodes[1].Storage, secondariesWithoutFailover[0].GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[1].Address, secondariesWithoutFailover[0].GetAddress())

	// shard with failover should have promoted a secondary to primary and demoted the primary to a secondary
	require.Equal(t, virtualStorages[0].Nodes[1].Storage, primary.GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[1].Address, primary.GetAddress())
	require.Len(t, secondaries, 1)
	require.Equal(t, virtualStorages[0].Nodes[0].Storage, secondaries[0].GetStorage())
	require.Equal(t, virtualStorages[0].Nodes[0].Address, secondaries[0].GetAddress())

	healthSrv1.SetServingStatus("", grpc_health_v1.HealthCheckResponse_UNKNOWN)
	nm.checkShards()

	_, err = nm.GetShard("virtual-storage-0")
	require.Error(t, err, "should return error since no nodes are healthy")
}

func TestMgr_GetSyncedNode(t *testing.T) {
	var sockets [4]string
	var srvs [4]*grpc.Server
	var healthSrvs [4]*health.Server
	for i := 0; i < 4; i++ {
		sockets[i] = testhelper.GetTemporaryGitalySocketFileName()
		srvs[i], healthSrvs[i] = testhelper.NewServerWithHealth(t, sockets[i])
		defer srvs[i].Stop()
	}

	vs0Primary := "unix://" + sockets[0]

	virtualStorages := []*config.VirtualStorage{
		{
			Name: "virtual-storage-0",
			Nodes: []*models.Node{
				{
					Storage:        "gitaly-0",
					Address:        vs0Primary,
					DefaultPrimary: true,
				},
				{
					Storage: "gitaly-1",
					Address: "unix://" + sockets[1],
				},
			},
		},
		{
			// second virtual storage needed to check there is no intersections between two even with same storage names
			Name: "virtual-storage-1",
			Nodes: []*models.Node{
				{
					// same storage name as in other virtual storage is used intentionally
					Storage:        "gitaly-1",
					Address:        "unix://" + sockets[2],
					DefaultPrimary: true,
				},
				{
					Storage: "gitaly-2",
					Address: "unix://" + sockets[3],
				},
			},
		},
	}

	conf := config.Config{
		VirtualStorages: virtualStorages,
		Failover:        config.Failover{Enabled: true},
	}

	mockHistogram := promtest.NewMockHistogramVec()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ackEvent := func(ds datastore.Datastore, job datastore.ReplicationJob, state datastore.JobState) datastore.ReplicationEvent {
		event := datastore.ReplicationEvent{Job: job}

		eevent, err := ds.Enqueue(ctx, event)
		require.NoError(t, err)

		devents, err := ds.Dequeue(ctx, eevent.Job.VirtualStorage, eevent.Job.TargetNodeStorage, 100500)
		require.NoError(t, err)
		require.Len(t, devents, 1)

		acks, err := ds.Acknowledge(ctx, state, []uint64{devents[0].ID})
		require.NoError(t, err)
		require.Equal(t, []uint64{devents[0].ID}, acks)
		return devents[0]
	}

	verify := func(scenario func(t *testing.T, nm Manager, ds datastore.Datastore)) func(*testing.T) {
		ds := datastore.Datastore{ReplicationEventQueue: datastore.NewMemoryReplicationEventQueue(conf)}

		nm, err := NewManager(testhelper.DiscardTestEntry(t), conf, nil, ds, mockHistogram)
		require.NoError(t, err)

		nm.Start(time.Duration(0), time.Hour)

		return func(t *testing.T) { scenario(t, nm, ds) }
	}

	t.Run("unknown virtual storage", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		_, err := nm.GetSyncedNode(ctx, "virtual-storage-unknown", "")
		require.True(t, errors.Is(err, ErrVirtualStorageNotExist))
	}))

	t.Run("no replication events", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		node, err := nm.GetSyncedNode(ctx, "virtual-storage-0", "no/matter")
		require.NoError(t, err)
		require.Contains(t, []string{vs0Primary, "unix://" + sockets[1]}, node.GetAddress())
	}))

	t.Run("last replication event is in 'ready'", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		_, err := ds.Enqueue(ctx, datastore.ReplicationEvent{
			Job: datastore.ReplicationJob{
				RelativePath:      "path/1",
				TargetNodeStorage: "gitaly-1",
				SourceNodeStorage: "gitaly-0",
				VirtualStorage:    "virtual-storage-0",
			},
		})
		require.NoError(t, err)

		node, err := nm.GetSyncedNode(ctx, "virtual-storage-0", "path/1")
		require.NoError(t, err)
		require.Equal(t, vs0Primary, node.GetAddress())
	}))

	t.Run("last replication event is in 'in_progress'", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		vs0Event, err := ds.Enqueue(ctx, datastore.ReplicationEvent{
			Job: datastore.ReplicationJob{
				RelativePath:      "path/1",
				TargetNodeStorage: "gitaly-1",
				SourceNodeStorage: "gitaly-0",
				VirtualStorage:    "virtual-storage-0",
			},
		})
		require.NoError(t, err)

		vs0Events, err := ds.Dequeue(ctx, vs0Event.Job.VirtualStorage, vs0Event.Job.TargetNodeStorage, 100500)
		require.NoError(t, err)
		require.Len(t, vs0Events, 1)

		node, err := nm.GetSyncedNode(ctx, "virtual-storage-0", "path/1")
		require.NoError(t, err)
		require.Equal(t, vs0Primary, node.GetAddress())
	}))

	t.Run("last replication event is in 'failed'", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		vs0Event := ackEvent(ds, datastore.ReplicationJob{
			RelativePath:      "path/1",
			TargetNodeStorage: "gitaly-1",
			SourceNodeStorage: "gitaly-0",
			VirtualStorage:    "virtual-storage-0",
		}, datastore.JobStateFailed)

		node, err := nm.GetSyncedNode(ctx, vs0Event.Job.VirtualStorage, vs0Event.Job.RelativePath)
		require.NoError(t, err)
		require.Equal(t, vs0Primary, node.GetAddress())
	}))

	t.Run("multiple replication events for same virtual, last is in 'ready'", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		vsEvent0 := ackEvent(ds, datastore.ReplicationJob{
			RelativePath:      "path/1",
			TargetNodeStorage: "gitaly-1",
			SourceNodeStorage: "gitaly-0",
			VirtualStorage:    "virtual-storage-0",
		}, datastore.JobStateCompleted)

		vsEvent1, err := ds.Enqueue(ctx, datastore.ReplicationEvent{
			Job: datastore.ReplicationJob{
				RelativePath:      vsEvent0.Job.RelativePath,
				TargetNodeStorage: vsEvent0.Job.TargetNodeStorage,
				SourceNodeStorage: vsEvent0.Job.SourceNodeStorage,
				VirtualStorage:    vsEvent0.Job.VirtualStorage,
			},
		})
		require.NoError(t, err)

		node, err := nm.GetSyncedNode(ctx, vsEvent1.Job.VirtualStorage, vsEvent1.Job.RelativePath)
		require.NoError(t, err)
		require.Equal(t, vs0Primary, node.GetAddress())
	}))

	t.Run("same repo path for different virtual storages", verify(func(t *testing.T, nm Manager, ds datastore.Datastore) {
		vs0Event := ackEvent(ds, datastore.ReplicationJob{
			RelativePath:      "path/1",
			TargetNodeStorage: "gitaly-1",
			SourceNodeStorage: "gitaly-0",
			VirtualStorage:    "virtual-storage-0",
		}, datastore.JobStateDead)

		ackEvent(ds, datastore.ReplicationJob{
			RelativePath:      "path/1",
			TargetNodeStorage: "gitaly-2",
			SourceNodeStorage: "gitaly-1",
			VirtualStorage:    "virtual-storage-1",
		}, datastore.JobStateCompleted)

		node, err := nm.GetSyncedNode(ctx, vs0Event.Job.VirtualStorage, vs0Event.Job.RelativePath)
		require.NoError(t, err)
		require.Equal(t, vs0Primary, node.GetAddress())
	}))
}
