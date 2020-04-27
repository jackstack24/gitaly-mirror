package praefect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/client"
	gconfig "gitlab.com/gitlab-org/gitaly/internal/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/models"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
)

func TestInfoService_RepositoryReplicas(t *testing.T) {
	conf := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			{
				Name: "default",
				Nodes: []*models.Node{
					{
						Storage:        "praefect-internal-1",
						DefaultPrimary: true,
					},
					{
						Storage: "praefect-internal-2",
					},
					{
						Storage: "praefect-internal-3",
					},
				},
			},
		},
	}

	testRepo, _, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	cc, _, cleanup := runPraefectServerWithGitaly(t, conf)
	defer cleanup()

	defer func(storages []gconfig.Storage) {
		gconfig.Config.Storages = storages
	}(gconfig.Config.Storages)

	tempDir, cleanupTempDir := testhelper.TempDir(t)
	defer cleanupTempDir()

	for _, node := range conf.VirtualStorages[0].Nodes {
		storagePath := filepath.Join(tempDir, node.Storage)
		require.NoError(t, os.MkdirAll(storagePath, 0755))
		gconfig.Config.Storages = append(gconfig.Config.Storages, gconfig.Storage{
			Name: node.Storage,
			Path: storagePath,
		})
	}

	testRepoPrimary, _, cleanup := cloneRepoAtStorage(t, testRepo, conf.VirtualStorages[0].Nodes[0].Storage)
	defer cleanup()

	_, _, cleanup = cloneRepoAtStorage(t, testRepo, conf.VirtualStorages[0].Nodes[1].Storage)
	defer cleanup()
	_, testRepoSecondary2Path, cleanup := cloneRepoAtStorage(t, testRepo, conf.VirtualStorages[0].Nodes[2].Storage)
	defer cleanup()

	// create a commit in the second replica so we can check that its checksum is different than the primary
	testhelper.CreateCommit(t, testRepoSecondary2Path, "master", nil)

	praefectInfoServiceClient := gitalypb.NewPraefectInfoServiceClient(cc)

	ctx, cancel := testhelper.Context()
	defer cancel()

	// get the checksum of the primary
	primaryGitaly, err := client.Dial(conf.VirtualStorages[0].Nodes[0].Address, nil)
	require.NoError(t, err)
	defer primaryGitaly.Close()

	primaryGitalyRepoClient := gitalypb.NewRepositoryServiceClient(primaryGitaly)
	checksum, err := primaryGitalyRepoClient.CalculateChecksum(ctx, &gitalypb.CalculateChecksumRequest{
		Repository: &gitalypb.Repository{
			StorageName:  conf.VirtualStorages[0].Name,
			RelativePath: testRepoPrimary.GetRelativePath(),
		},
	})
	require.NoError(t, err)

	resp, err := praefectInfoServiceClient.RepositoryReplicas(ctx, &gitalypb.RepositoryReplicasRequest{
		Repository: &gitalypb.Repository{
			StorageName:  conf.VirtualStorages[0].Name,
			RelativePath: testRepoPrimary.GetRelativePath(),
		},
	})

	require.NoError(t, err)

	require.Equal(t, checksum.Checksum, resp.Primary.Checksum)
	var checked []string
	for _, secondary := range resp.GetReplicas() {
		switch storage := secondary.GetRepository().GetStorageName(); storage {
		case conf.VirtualStorages[0].Nodes[1].Storage:
			require.Equal(t, checksum.Checksum, secondary.Checksum)
			checked = append(checked, storage)
		case conf.VirtualStorages[0].Nodes[2].Storage:
			require.NotEqual(t, checksum.Checksum, secondary.Checksum, "should not be equal since we added a commit")
			checked = append(checked, storage)
		default:
			require.FailNow(t, "unexpected storage: %q", storage)
		}
	}
	require.ElementsMatch(t, []string{conf.VirtualStorages[0].Nodes[1].Storage, conf.VirtualStorages[0].Nodes[2].Storage}, checked)
}
