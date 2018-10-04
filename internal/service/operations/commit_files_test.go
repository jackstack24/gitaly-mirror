	"fmt"
		desc            string
		repo            *pb.Repository
		repoPath        string
		branchName      string
		repoCreated     bool
		branchCreated   bool
		executeFilemode bool
		{
			desc:            "create executable file",
			repo:            testRepo,
			repoPath:        testRepoPath,
			branchName:      "feature-executable",
			repoCreated:     false,
			branchCreated:   true,
			executeFilemode: true,
		},
			actionsRequest4 := chmodFileHeaderRequest(filePath, tc.executeFilemode)
			require.NoError(t, stream.Send(actionsRequest4))

			commitInfo := testhelper.MustRunCommand(t, nil, "git", "-C", tc.repoPath, "show", headCommit.GetId())
			expectedFilemode := "100644"
			if tc.executeFilemode {
				expectedFilemode = "100755"
			}
			require.Contains(t, string(commitInfo), fmt.Sprint("new file mode ", expectedFilemode))
		{
			desc: "file doesn't exists",
			requests: []*pb.UserCommitFilesRequest{
				headerRequest(testRepo, user, "feature", commitFilesMessage, nil, nil),
				chmodFileHeaderRequest("documents/story.txt", true),
			},
			indexError: "A file with this name doesn't exist",
		},
func chmodFileHeaderRequest(filePath string, executeFilemode bool) *pb.UserCommitFilesRequest {
	return actionRequest(&pb.UserCommitFilesAction{
		UserCommitFilesActionPayload: &pb.UserCommitFilesAction_Header{
			Header: &pb.UserCommitFilesActionHeader{
				Action:          pb.UserCommitFilesActionHeader_CHMOD,
				FilePath:        []byte(filePath),
				ExecuteFilemode: executeFilemode,
			},
		},
	})
}
