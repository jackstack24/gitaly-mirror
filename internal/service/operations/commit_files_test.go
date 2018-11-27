	"strconv"
func TestSuccessfulUserCommitFilesRequestMove(t *testing.T) {
	server, serverSocketPath := runFullServer(t)
	defer server.Stop()

	client, conn := operations.NewOperationClient(t, serverSocketPath)
	defer conn.Close()

	ctxOuter, cancel := testhelper.Context()
	defer cancel()

	branchName := "master"
	previousFilePath := "README"
	filePath := "NEWREADME"
	authorName := []byte("Jane Doe")
	authorEmail := []byte("janedoe@gitlab.com")

	for i, tc := range []struct {
		content string
		infer   bool
	}{
		{content: "", infer: false},
		{content: "foo", infer: false},
		{content: "", infer: true},
		{content: "foo", infer: true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			testRepo, testRepoPath, cleanupFn := testhelper.NewTestRepo(t)
			defer cleanupFn()

			origFileContent := testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "show", branchName+":"+previousFilePath)
			md := testhelper.GitalyServersMetadata(t, serverSocketPath)
			ctx := metadata.NewOutgoingContext(ctxOuter, md)
			headerRequest := headerRequest(testRepo, user, branchName, commitFilesMessage, authorName, authorEmail)
			actionsRequest1 := moveFileHeaderRequest(previousFilePath, filePath, tc.infer)

			stream, err := client.UserCommitFiles(ctx)
			require.NoError(t, err)
			require.NoError(t, stream.Send(headerRequest))
			require.NoError(t, stream.Send(actionsRequest1))

			if len(tc.content) > 0 {
				actionsRequest2 := actionContentRequest(tc.content)
				require.NoError(t, stream.Send(actionsRequest2))
			}

			r, err := stream.CloseAndRecv()
			require.NoError(t, err)

			update := r.GetBranchUpdate()
			require.NotNil(t, update)

			fileContent := testhelper.MustRunCommand(t, nil, "git", "-C", testRepoPath, "show", update.CommitId+":"+filePath)

			if tc.infer {
				require.Equal(t, string(origFileContent), string(fileContent))
			} else {
				require.Equal(t, tc.content, string(fileContent))
			}
		})
	}
}

func moveFileHeaderRequest(previousPath, filePath string, infer bool) *gitalypb.UserCommitFilesRequest {
	return actionRequest(&gitalypb.UserCommitFilesAction{
		UserCommitFilesActionPayload: &gitalypb.UserCommitFilesAction_Header{
			Header: &gitalypb.UserCommitFilesActionHeader{
				Action:       gitalypb.UserCommitFilesActionHeader_MOVE,
				FilePath:     []byte(filePath),
				PreviousPath: []byte(previousPath),
				InferContent: infer,
			},
		},
	})
}
