package praefect

import (
	"context"
	"net"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitalyauth "gitlab.com/gitlab-org/gitaly/auth"
	"gitlab.com/gitlab-org/gitaly/internal/config/auth"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/mock"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/models"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/nodes"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/protoregistry"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/transactions"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper/promtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestAuthFailures(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	testCases := []struct {
		desc string
		opts []grpc.DialOption
		code codes.Code
	}{
		{
			desc: "no auth",
			opts: nil,
			code: codes.Unauthenticated,
		},
		{
			desc: "invalid auth",
			opts: []grpc.DialOption{grpc.WithPerRPCCredentials(brokenAuth{})},
			code: codes.Unauthenticated,
		},
		{
			desc: "wrong secret new auth",
			opts: []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2("foobar"))},
			code: codes.PermissionDenied,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			srv, serverSocketPath, cleanup := runServer(t, "quxbaz", true)
			defer srv.Shutdown(ctx)
			defer cleanup()

			connOpts := append(tc.opts, grpc.WithInsecure())
			conn, err := dial(serverSocketPath, connOpts)
			require.NoError(t, err, tc.desc)
			defer conn.Close()

			cli := mock.NewSimpleServiceClient(conn)

			_, err = cli.ServerAccessor(ctx, &mock.SimpleRequest{
				Value: 1,
			})

			testhelper.RequireGrpcError(t, err, tc.code)
		})
	}
}

func TestAuthSuccess(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	token := "foobar"

	testCases := []struct {
		desc     string
		opts     []grpc.DialOption
		required bool
		token    string
	}{
		{desc: "no auth, not required"},
		{
			desc:  "v2 correct auth, not required",
			opts:  []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2(token))},
			token: token,
		},
		{
			desc:  "v2 incorrect auth, not required",
			opts:  []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2("incorrect"))},
			token: token,
		},
		{
			desc:     "v2 correct auth, required",
			opts:     []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2(token))},
			token:    token,
			required: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			srv, serverSocketPath, cleanup := runServer(t, tc.token, tc.required)
			defer srv.Shutdown(ctx)
			defer cleanup()

			connOpts := append(tc.opts, grpc.WithInsecure())
			conn, err := dial(serverSocketPath, connOpts)
			require.NoError(t, err, tc.desc)
			defer conn.Close()

			cli := mock.NewSimpleServiceClient(conn)

			_, err = cli.ServerAccessor(ctx, &mock.SimpleRequest{
				Value: 1,
			})

			assert.NoError(t, err, tc.desc)
		})
	}
}

type brokenAuth struct{}

func (brokenAuth) RequireTransportSecurity() bool { return false }
func (brokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer blablabla"}, nil
}

func dial(serverSocketPath string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(serverSocketPath, opts...)
}

func runServer(t *testing.T, token string, required bool) (*Server, string, func()) {
	backendToken := "abcxyz"
	mockServer := &mockSvc{
		serverAccessor: func(_ context.Context, req *mock.SimpleRequest) (*mock.SimpleResponse, error) {
			return &mock.SimpleResponse{
				Value: req.Value + 1,
			}, nil
		},
	}
	backend, cleanup := newMockDownstream(t, backendToken, mockServer)

	conf := config.Config{
		Auth: auth.Config{Token: token, Transitioning: !required},
		VirtualStorages: []*config.VirtualStorage{
			&config.VirtualStorage{
				Name: "praefect",
				Nodes: []*models.Node{
					&models.Node{
						Storage:        "praefect-internal-0",
						DefaultPrimary: true,
						Address:        backend,
						Token:          backendToken,
					},
				},
			},
		},
	}

	gz := proto.FileDescriptor("mock.proto")
	fd, err := protoregistry.ExtractFileDescriptor(gz)
	if err != nil {
		t.Fatal(err)
	}

	logEntry := testhelper.DiscardTestEntry(t)
	ds := datastore.Datastore{
		ReplicasDatastore:     datastore.NewInMemory(conf),
		ReplicationEventQueue: datastore.NewMemoryReplicationEventQueue(),
	}

	nodeMgr, err := nodes.NewManager(logEntry, conf, nil, promtest.NewMockHistogramVec())
	require.NoError(t, err)

	txMgr := transactions.NewManager()

	registry := protoregistry.New()
	require.NoError(t, registry.RegisterFiles(fd))

	coordinator := NewCoordinator(logEntry, ds, nodeMgr, txMgr, conf, registry)

	srv := NewServer(coordinator.StreamDirector, logEntry, registry, conf)

	serverSocketPath := testhelper.GetTemporaryGitalySocketFileName()

	listener, err := net.Listen("unix", serverSocketPath)
	require.NoError(t, err)
	go srv.Serve(listener, false)

	return srv, "unix://" + serverSocketPath, cleanup
}
