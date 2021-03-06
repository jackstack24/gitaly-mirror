package ref

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/gitlab-org/gitaly/internal/command"
	"gitlab.com/gitlab-org/gitaly/internal/git"
	"gitlab.com/gitlab-org/gitaly/internal/helper"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
)

// RefExists returns true if the given reference exists. The ref must start with the string `ref/`
func (server) RefExists(ctx context.Context, in *gitalypb.RefExistsRequest) (*gitalypb.RefExistsResponse, error) {
	ref := string(in.Ref)

	if !isValidRefName(ref) {
		return nil, helper.ErrInvalidArgument(fmt.Errorf("invalid refname"))
	}

	exists, err := refExists(ctx, in.Repository, ref)
	if err != nil {
		return nil, helper.ErrInternal(err)
	}

	return &gitalypb.RefExistsResponse{Value: exists}, nil
}

func refExists(ctx context.Context, repo *gitalypb.Repository, ref string) (bool, error) {
	cmd, err := git.SafeCmd(ctx, repo, nil, git.SubCmd{
		Name:  "show-ref",
		Flags: []git.Option{git.Flag{"--verify"}, git.Flag{"--quiet"}},
		Args:  []string{ref},
	})
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err == nil {
		// Exit code 0: the ref exists
		return true, nil
	}

	if code, ok := command.ExitStatus(err); ok && code == 1 {
		// Exit code 1: the ref does not exist
		return false, nil
	}

	// This will normally occur when exit code > 1
	return false, err
}

func isValidRefName(refName string) bool {
	return strings.HasPrefix(refName, "refs/")
}
