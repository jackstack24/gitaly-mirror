package objectpool

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/git"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
)

func TestLink(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	testRepo, _, cleanupFn := testhelper.NewTestRepo(t)
	defer cleanupFn()

	pool, poolCleanup := NewTestObjectPool(ctx, t, testRepo.GetStorageName())
	defer poolCleanup()

	require.NoError(t, pool.Remove(ctx), "make sure pool does not exist prior to creation")
	require.NoError(t, pool.Create(ctx, testRepo), "create pool")

	altPath, err := git.InfoAlternatesPath(testRepo)
	require.NoError(t, err)
	_, err = os.Stat(altPath)
	require.True(t, os.IsNotExist(err))

	require.NoError(t, pool.Link(ctx, testRepo))

	require.FileExists(t, altPath, "alternates file must exist after Link")

	content, err := ioutil.ReadFile(altPath)
	require.NoError(t, err)

	require.True(t, strings.HasPrefix(string(content), "../"), "expected %q to be relative path", content)

	require.NoError(t, pool.Link(ctx, testRepo))

	newContent, err := ioutil.ReadFile(altPath)
	require.NoError(t, err)

	require.Equal(t, content, newContent)

	require.False(t, testhelper.RemoteExists(t, pool.FullPath(), testRepo.GetGlRepository()), "pool remotes should not include %v", testRepo)
}

func TestLinkRemoveBitmap(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	testCases := []struct {
		desc           string
		poolHasBitmaps bool
	}{
		{
			desc:           "pool has bitmaps",
			poolHasBitmaps: true,
		},
		{
			desc:           "pool does not have bitmaps",
			poolHasBitmaps: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testRepo, testRepoPath, cleanupFn := testhelper.NewTestRepo(t)
			defer cleanupFn()

			pool, poolCleanup := NewTestObjectPool(ctx, t, testRepo.GetStorageName())
			defer poolCleanup()

			require.NoError(t, pool.Init(ctx))

			poolPath := pool.FullPath()
			testhelper.MustRunCommand(t, nil, "git", "-C", poolPath, "fetch", testRepoPath, "+refs/*:refs/*")

			testhelper.MustRunCommand(t, nil, "git", "-C", poolPath, "repack", "-adb")
			if tc.poolHasBitmaps {
				require.Len(t, listBitmaps(t, pool.FullPath()), 1, "pool bitmaps before")
			} else {
				removeBitmaps(t, pool.FullPath())
				require.Len(t, listBitmaps(t, pool.FullPath()), 0, "pool bitmaps before")
			}

			testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "repack", "-adb")
			require.Len(t, listBitmaps(t, testRepoPath), 1, "member bitmaps before")

			refsBefore := testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "for-each-ref")

			require.NoError(t, pool.Link(ctx, testRepo))

			if tc.poolHasBitmaps {
				require.Len(t, listBitmaps(t, pool.FullPath()), 1, "pool bitmaps after")
			} else {
				require.Len(t, listBitmaps(t, pool.FullPath()), 0, "pool bitmaps after")
			}
			require.Len(t, listBitmaps(t, testRepoPath), 0, "member bitmaps after")

			testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "fsck")

			refsAfter := testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "for-each-ref")
			require.Equal(t, refsBefore, refsAfter, "compare member refs before/after link")
		})
	}
}

func removeBitmaps(t *testing.T, repoPath string) {
	packDir := filepath.Join(repoPath, "objects", "pack")

	entries, err := ioutil.ReadDir(packDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".bitmap") {
			bitmap := filepath.Join(packDir, entry.Name())
			require.NoError(t, os.Remove(bitmap))
		}
	}
}

func listBitmaps(t *testing.T, repoPath string) []string {
	entries, err := ioutil.ReadDir(filepath.Join(repoPath, "objects/pack"))
	require.NoError(t, err)

	var bitmaps []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".bitmap") {
			bitmaps = append(bitmaps, entry.Name())
		}
	}

	return bitmaps
}

func TestUnlink(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	testRepo, _, cleanupFn := testhelper.NewTestRepo(t)
	defer cleanupFn()

	pool, poolCleanup := NewTestObjectPool(ctx, t, testRepo.GetStorageName())
	defer poolCleanup()

	require.Error(t, pool.Unlink(ctx, testRepo), "removing a non-existing pool should be an error")

	require.NoError(t, pool.Create(ctx, testRepo), "create pool")
	require.NoError(t, pool.Link(ctx, testRepo), "link test repo to pool")

	require.False(t, testhelper.RemoteExists(t, pool.FullPath(), testRepo.GetGlRepository()), "pool remotes should include %v", testRepo)

	require.NoError(t, pool.Unlink(ctx, testRepo), "unlink repo")
	require.False(t, testhelper.RemoteExists(t, pool.FullPath(), testRepo.GetGlRepository()), "pool remotes should no longer include %v", testRepo)
}

func TestLinkAbsoluteLinkExists(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	testRepo, testRepoPath, cleanupFn := testhelper.NewTestRepo(t)
	defer cleanupFn()

	pool, poolCleanup := NewTestObjectPool(ctx, t, testRepo.GetStorageName())
	defer poolCleanup()

	require.NoError(t, pool.Remove(ctx), "make sure pool does not exist prior to creation")
	require.NoError(t, pool.Create(ctx, testRepo), "create pool")

	altPath, err := git.InfoAlternatesPath(testRepo)
	require.NoError(t, err)

	fullPath := filepath.Join(pool.FullPath(), "objects")

	require.NoError(t, ioutil.WriteFile(altPath, []byte(fullPath), 0644))

	require.NoError(t, pool.Link(ctx, testRepo), "we expect this call to change the absolute link to a relative link")

	require.FileExists(t, altPath, "alternates file must exist after Link")

	content, err := ioutil.ReadFile(altPath)
	require.NoError(t, err)

	require.False(t, filepath.IsAbs(string(content)), "expected %q to be relative path", content)

	testRepoObjectsPath := filepath.Join(testRepoPath, "objects")
	require.Equal(t, fullPath, filepath.Join(testRepoObjectsPath, string(content)), "the content of the alternates file should be the relative version of the absolute pat")

	require.True(t, testhelper.RemoteExists(t, pool.FullPath(), "origin"), "pool remotes should include %v", testRepo)
}
