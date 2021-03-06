package objectpool

import (
	"context"

	"gitlab.com/gitlab-org/gitaly/internal/git"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) ReduplicateRepository(ctx context.Context, req *gitalypb.ReduplicateRepositoryRequest) (*gitalypb.ReduplicateRepositoryResponse, error) {
	if req.GetRepository() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ReduplicateRepository: no repository")
	}

	cmd, err := git.SafeCmd(ctx, req.GetRepository(), nil, git.SubCmd{
		Name:  "repack",
		Flags: []git.Option{git.Flag{"--quiet"}, git.Flag{"-a"}},
	})
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return &gitalypb.ReduplicateRepositoryResponse{}, nil
}
