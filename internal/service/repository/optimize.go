package repository

import (
	"context"

	"gitlab.com/gitlab-org/gitaly/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
)

func (s *server) optimizeRepository(ctx context.Context, repository *gitalypb.Repository) error {
	profile, err := stats.GetProfile(ctx, repository)
	if err != nil {
		return err
	}

	if !profile.HasBitmap() {
		_, err := s.RepackFull(ctx, &gitalypb.RepackFullRequest{Repository: repository, CreateBitmap: true})
		if err != nil {
			return err
		}
	}

	return nil
}
