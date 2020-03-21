package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
)

func getNewestPackfileModtime(t *testing.T, repoPath string) time.Time {
	packFiles, err := filepath.Glob(filepath.Join(repoPath, "objects", "pack", "*.pack"))
	require.NoError(t, err)
	if len(packFiles) == 0 {
		t.Error("no packfiles exist")
	}

	var newestPackfileModtime time.Time

	for _, packFile := range packFiles {
		info, err := os.Stat(packFile)
		require.NoError(t, err)
		if info.ModTime().After(newestPackfileModtime) {
			newestPackfileModtime = info.ModTime()
		}
	}

	return newestPackfileModtime
}

func TestOptimizeRepository(t *testing.T) {
	testRepo, testRepoPath, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "repack", "-A", "-b")

	ctx, cancel := testhelper.Context()
	defer cancel()

	profile, err := stats.GetProfile(ctx, testRepo)
	require.NoError(t, err)

	require.True(t, profile.HasBitmap())

	// get timestamp of latest packfile
	newestsPackfileTime := getNewestPackfileModtime(t, testRepoPath)

	testhelper.CreateCommit(t, testRepoPath, "master", nil)

	repoServer := &server{}

	require.NoError(t, repoServer.optimizeRepository(ctx, testRepo))
	require.Equal(t, getNewestPackfileModtime(t, testRepoPath), newestsPackfileTime, "there should not have been a new packfile created")

	testRepo, testRepoPath, cleanup = testhelper.InitBareRepo(t)

	blobs := 10
	blobIDs := testhelper.WriteBlobs(t, testRepoPath, blobs)

	for _, blobID := range blobIDs {
		commitID := testhelper.CommitBlobWithName(t, testRepoPath, blobID, blobID, "adding another blob....")
		testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "update-ref", "refs/heads/"+blobID, commitID)
	}

	bitmaps, err := filepath.Glob(filepath.Join(testRepoPath, "objects", "pack", "*.bitmap"))
	require.NoError(t, err)
	require.Empty(t, bitmaps)

	// optimize repository on a repository without a bitmap should call repack full
	require.NoError(t, repoServer.optimizeRepository(ctx, testRepo))

	bitmaps, err = filepath.Glob(filepath.Join(testRepoPath, "objects", "pack", "*.bitmap"))
	require.NoError(t, err)
	require.NotEmpty(t, bitmaps)
}
