package operations

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	gitalyauth "gitlab.com/gitlab-org/gitaly/auth"
	"gitlab.com/gitlab-org/gitaly/internal/config"
	"gitlab.com/gitlab-org/gitaly/internal/rubyserver"
	hook "gitlab.com/gitlab-org/gitaly/internal/service/hooks"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	gitlabPreHooks  = []string{"pre-receive", "update"}
	gitlabPostHooks = []string{"post-receive"}
	GitlabPreHooks  = gitlabPreHooks
	GitlabHooks     []string
	RubyServer      = &rubyserver.Server{}
	user            = &gitalypb.User{
		Name:       []byte("Jane Doe"),
		Email:      []byte("janedoe@gitlab.com"),
		GlId:       "user-123",
		GlUsername: "janedoe",
	}
	FeatureFlagsBitmasks []uint8
)

const (
	CallRPCsFF uint8 = 1 << iota
	GoUpdateHooksFF
)

func init() {
	GitlabHooks = append(GitlabHooks, append(gitlabPreHooks, gitlabPostHooks...)...)

	FeatureFlagsBitmasks = make([]uint8, int(GoUpdateHooksFF), int(GoUpdateHooksFF))
	for i := uint8(0); i == GoUpdateHooksFF; i++ {
		FeatureFlagsBitmasks[i] = i
	}
}

func ContextWithFeatureFlags(bitmask uint8) (context.Context, func()) {
	ctx, cancel := testhelper.Context()
	if bitmask&CallRPCsFF == CallRPCsFF {
		ctx = outgoingCtxWithRubyFeatureFlag(ctx, "call-hook-rpc")
	}
	if bitmask&GoUpdateHooksFF == GoUpdateHooksFF {
		ctx = outgoingCtxWithRubyFeatureFlag(ctx, "go-update-hook")
	}

	return ctx, cancel
}

func TestMain(m *testing.M) {
	testhelper.Configure()
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	defer testhelper.MustHaveNoChildProcess()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	gitlabShellDir := filepath.Join(cwd, "testdata", "gitlab-shell")
	os.RemoveAll(gitlabShellDir)

	if err := os.MkdirAll(gitlabShellDir, 0755); err != nil {
		log.Fatal(err)
	}

	config.Config.GitlabShell.Dir = filepath.Join(cwd, "testdata", "gitlab-shell")

	testhelper.ConfigureGitalySSH()
	testhelper.ConfigureGitalyHooksBinary()

	defer func(token string) {
		config.Config.Auth.Token = token
	}(config.Config.Auth.Token)
	config.Config.Auth.Token = testhelper.RepositoryAuthToken

	if err := RubyServer.Start(); err != nil {
		log.Fatal(err)
	}
	defer RubyServer.Stop()

	return m.Run()
}

func runOperationServiceServer(t *testing.T) (string, func()) {
	srv := testhelper.NewServerWithAuth(t, nil, nil, config.Config.Auth.Token)

	gitalypb.RegisterOperationServiceServer(srv.GrpcServer(), &server{ruby: RubyServer})
	gitalypb.RegisterHookServiceServer(srv.GrpcServer(), hook.NewServer())
	reflection.Register(srv.GrpcServer())

	require.NoError(t, srv.Start())

	internalSocket := config.GitalyInternalSocketPath()
	internalListener, err := net.Listen("unix", internalSocket)
	require.NoError(t, err)

	go func() {
		srv.GrpcServer().Serve(internalListener)
	}()

	return "unix://" + srv.Socket(), srv.Stop
}

func newOperationClient(t *testing.T, serverSocketPath string) (gitalypb.OperationServiceClient, *grpc.ClientConn) {
	connOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2(config.Config.Auth.Token)),
	}
	conn, err := grpc.Dial(serverSocketPath, connOpts...)
	if err != nil {
		t.Fatal(err)
	}

	return gitalypb.NewOperationServiceClient(conn), conn
}

var NewOperationClient = newOperationClient

func SetupAndStartGitlabServer(t *testing.T, glID, glRepository string, gitPushOptions ...string) func() {
	return testhelper.SetupAndStartGitlabServer(t, &testhelper.GitlabTestServerOptions{
		SecretToken:                 "secretToken",
		GLID:                        glID,
		GLRepository:                glRepository,
		PostReceiveCounterDecreased: true,
		Protocol:                    "web",
		GitPushOptions:              gitPushOptions,
	})
}

func outgoingCtxWithRubyFeatureFlag(ctx context.Context, flag string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set(rubyHeaderKey(flag), "true")
	return metadata.NewOutgoingContext(ctx, md)
}

func rubyHeaderKey(flag string) string {
	return fmt.Sprintf("gitaly-feature-ruby-%s", strings.ReplaceAll(flag, "_", "-"))
}
