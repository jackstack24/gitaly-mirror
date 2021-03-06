package hook

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
)

// printAllScript is a bash script that prints out stdin, the arguments,
// and the environment variables in the following format:
// stdin:old new ref0
// args: arg1 arg2
// env: VAR1=VAL1 VAR2=VAL2
// NOTE: this script only prints one line of stdin
var printAllScript = []byte(`#!/bin/bash
read stdin
echo stdin:$stdin
echo args:$@
echo env: $(printenv)`)

// printStdinScript prints stdin line by line
var printStdinScript = []byte(`#!/bin/bash
while read line
do
  echo "$line"
done
`)

// failScript prints the name of the command and exits with exit code 1
var failScript = []byte(`#!/bin/bash
echo "$0" >&2
exit 1`)

// successScript prints the name of the command and exits with exit code 0
var successScript = []byte(`#!/bin/bash
echo "$0"
exit 0`)

func TestCustomHooksSuccess(t *testing.T) {
	_, testRepoPath, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	testCases := []struct {
		hookName string
		stdin    string
		args     []string
		env      []string
		hookDir  string
	}{
		{
			hookName: "pre-receive",
			stdin:    "old new ref0",
			args:     nil,
			env:      []string{"GL_ID=user-123", "GL_USERNAME=username123", "GL_PROTOCOL=ssh", "GL_REPOSITORY=repo1"},
		},
		{
			hookName: "update",
			stdin:    "",
			args:     []string{"old", "new", "ref0"},
			env:      []string{"GL_ID=user-123", "GL_USERNAME=username123", "GL_PROTOCOL=ssh", "GL_REPOSITORY=repo1"},
		},
		{
			hookName: "post-receive",
			stdin:    "old new ref1",
			args:     nil,
			env:      []string{"GL_ID=user-123", "GL_USERNAME=username123", "GL_PROTOCOL=ssh", "GL_REPOSITORY=repo1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.hookName, func(t *testing.T) {
			globalCustomHooksDir, cleanupGlobalDir := testhelper.TempDir(t)
			defer cleanupGlobalDir()

			// hook is in project custom hook directory <repository>.git/custom_hooks/<hook_name>
			hookDir := filepath.Join(testRepoPath, "custom_hooks")
			callAndVerifyHooks(t, testRepoPath, tc.hookName, globalCustomHooksDir, hookDir, tc.stdin, tc.args, tc.env)

			// hook is in project custom hooks directory <repository>.git/custom_hooks/<hook_name>.d/*
			hookDir = filepath.Join(testRepoPath, "custom_hooks", fmt.Sprintf("%s.d", tc.hookName))
			callAndVerifyHooks(t, testRepoPath, tc.hookName, globalCustomHooksDir, hookDir, tc.stdin, tc.args, tc.env)

			// hook is in global custom hooks directory <global_custom_hooks_dir>/<hook_name>.d/*
			hookDir = filepath.Join(globalCustomHooksDir, fmt.Sprintf("%s.d", tc.hookName))
			callAndVerifyHooks(t, testRepoPath, tc.hookName, globalCustomHooksDir, hookDir, tc.stdin, tc.args, tc.env)
		})
	}
}

func TestCustomHookPartialFailure(t *testing.T) {
	_, testRepoPath, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	globalCustomHooksDir, cleanup := testhelper.TempDir(t)
	defer cleanup()

	ctx, cancel := testhelper.Context()
	defer cancel()

	testCases := []struct {
		hook                string
		projectHookSucceeds bool
		globalHookSucceeds  bool
	}{
		{
			hook:                "pre-receive",
			projectHookSucceeds: true,
			globalHookSucceeds:  false,
		},
		{
			hook:                "post-receive",
			projectHookSucceeds: false,
			globalHookSucceeds:  true,
		},
		{
			hook:                "update",
			projectHookSucceeds: false,
			globalHookSucceeds:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.hook, func(t *testing.T) {
			projectHookScript := successScript
			if !tc.projectHookSucceeds {
				projectHookScript = failScript
			}
			projectHookPath := filepath.Join(testRepoPath, "custom_hooks")
			cleanup := writeCustomHook(t, tc.hook, projectHookPath, projectHookScript)
			defer cleanup()

			globalHookScript := successScript
			if !tc.globalHookSucceeds {
				globalHookScript = failScript
			}
			globalHookPath := filepath.Join(globalCustomHooksDir, fmt.Sprintf("%s.d", tc.hook))
			cleanup = writeCustomHook(t, tc.hook, globalHookPath, globalHookScript)
			defer cleanup()

			caller, err := newCustomHooksExecutor(testRepoPath, globalCustomHooksDir, tc.hook)
			require.NoError(t, err)

			var stdout, stderr bytes.Buffer
			require.Error(t, caller(ctx, nil, nil, &bytes.Buffer{}, &stdout, &stderr))

			if tc.projectHookSucceeds && tc.globalHookSucceeds {
				require.Equal(t, filepath.Join(projectHookPath, tc.hook), text.ChompBytes(stdout.Bytes()))
				require.Equal(t, filepath.Join(globalHookPath, tc.hook), text.ChompBytes(stdout.Bytes()))
			} else if tc.projectHookSucceeds && !tc.globalHookSucceeds {
				require.Equal(t, filepath.Join(projectHookPath, tc.hook), text.ChompBytes(stdout.Bytes()))
				require.Equal(t, filepath.Join(globalHookPath, tc.hook), text.ChompBytes(stderr.Bytes()))
			} else {
				require.Equal(t, filepath.Join(projectHookPath, tc.hook), text.ChompBytes(stderr.Bytes()))
			}
		})
	}
}

func TestCustomHooksMultipleHooks(t *testing.T) {
	_, testRepoPath, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	globalCustomHooksDir, cleanup := testhelper.TempDir(t)
	defer cleanup()

	ctx, cancel := testhelper.Context()
	defer cancel()

	var expectedExecutedScripts []string

	projectUpdateHooks := 9
	projectHooksPath := filepath.Join(testRepoPath, "custom_hooks", "update.d")

	for i := 0; i < projectUpdateHooks; i++ {
		fileName := fmt.Sprintf("update_%d", i)
		writeCustomHook(t, fileName, projectHooksPath, successScript)
		expectedExecutedScripts = append(expectedExecutedScripts, filepath.Join(projectHooksPath, fileName))
	}

	globalUpdateHooks := 6
	globalHooksPath := filepath.Join(globalCustomHooksDir, "update.d")
	for i := 0; i < globalUpdateHooks; i++ {
		fileName := fmt.Sprintf("update_%d", i)
		writeCustomHook(t, fileName, globalHooksPath, successScript)
		expectedExecutedScripts = append(expectedExecutedScripts, filepath.Join(globalHooksPath, fileName))
	}

	hooksExecutor, err := newCustomHooksExecutor(testRepoPath, globalCustomHooksDir, "update")
	require.NoError(t, err)

	var stdout, stderr bytes.Buffer
	require.NoError(t, hooksExecutor(ctx, nil, nil, &bytes.Buffer{}, &stdout, &stderr))
	require.Empty(t, stderr.Bytes())

	outputScanner := bufio.NewScanner(&stdout)

	for _, expectedScript := range expectedExecutedScripts {
		require.True(t, outputScanner.Scan())
		require.Equal(t, expectedScript, outputScanner.Text())
	}
}

func TestMultilineStdin(t *testing.T) {
	_, testRepoPath, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	globalCustomHooksDir, cleanup := testhelper.TempDir(t)
	defer cleanup()

	ctx, cancel := testhelper.Context()
	defer cancel()

	projectHooksPath := filepath.Join(testRepoPath, "custom_hooks", "pre-receive.d")

	writeCustomHook(t, "pre-receive-script", projectHooksPath, printStdinScript)

	hooksExecutor, err := newCustomHooksExecutor(testRepoPath, globalCustomHooksDir, "pre-receive")
	require.NoError(t, err)

	changes := `old1 new1 ref1
old2 new2 ref2
old3 new3 ref3
`
	stdin := bytes.NewBufferString(changes)
	var stdout, stderr bytes.Buffer

	require.NoError(t, hooksExecutor(ctx, nil, nil, stdin, &stdout, &stderr))
	require.Equal(t, changes, stdout.String())
}

func TestMultipleScriptsStdin(t *testing.T) {
	_, testRepoPath, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	globalCustomHooksDir, cleanup := testhelper.TempDir(t)
	defer cleanup()

	ctx, cancel := testhelper.Context()
	defer cancel()

	projectUpdateHooks := 9
	projectHooksPath := filepath.Join(testRepoPath, "custom_hooks", "pre-receive.d")

	for i := 0; i < projectUpdateHooks; i++ {
		fileName := fmt.Sprintf("pre-receive_%d", i)
		writeCustomHook(t, fileName, projectHooksPath, printStdinScript)
	}

	hooksExecutor, err := newCustomHooksExecutor(testRepoPath, globalCustomHooksDir, "pre-receive")
	require.NoError(t, err)

	changes := "oldref11 newref00 ref123445"

	var stdout, stderr bytes.Buffer
	require.NoError(t, hooksExecutor(ctx, nil, nil, bytes.NewBufferString(changes+"\n"), &stdout, &stderr))
	require.Empty(t, stderr.Bytes())

	outputScanner := bufio.NewScanner(&stdout)

	for i := 0; i < projectUpdateHooks; i++ {
		require.True(t, outputScanner.Scan())
		require.Equal(t, changes, outputScanner.Text())
	}
}

func callAndVerifyHooks(t *testing.T, repoPath, hookName, globalHooksDir, hookDir, stdin string, args, env []string) {
	ctx, cancel := testhelper.Context()
	defer cancel()
	var stdout, stderr bytes.Buffer

	cleanup := writeCustomHook(t, hookName, hookDir, printAllScript)
	defer cleanup()

	callHooks, err := newCustomHooksExecutor(repoPath, globalHooksDir, hookName)
	require.NoError(t, err)

	require.NoError(t, callHooks(ctx, args, env, bytes.NewBufferString(stdin), &stdout, &stderr))
	require.Empty(t, stderr.Bytes())

	results := getCustomHookResults(&stdout)
	assert.Equal(t, stdin, results.stdin)
	assert.Equal(t, args, results.args)
	assert.Subset(t, results.env, env)
}

func getCustomHookResults(stdout *bytes.Buffer) customHookResults {
	lines := strings.SplitN(stdout.String(), "\n", 3)
	stdinLine := strings.SplitN(strings.TrimSpace(lines[0]), ":", 2)
	argsLine := strings.SplitN(strings.TrimSpace(lines[1]), ":", 2)
	envLine := strings.SplitN(strings.TrimSpace(lines[2]), ":", 2)

	var args, env []string
	if len(argsLine) == 2 && argsLine[1] != "" {
		args = strings.Split(argsLine[1], " ")
	}
	if len(envLine) == 2 && envLine[1] != "" {
		env = strings.Split(envLine[1], " ")
	}

	var stdin string
	if len(stdinLine) == 2 {
		stdin = stdinLine[1]
	}

	return customHookResults{
		stdin: stdin,
		args:  args,
		env:   env,
	}
}

type customHookResults struct {
	stdin string
	args  []string
	env   []string
}

func writeCustomHook(t *testing.T, hookName, dir string, content []byte) func() {
	require.NoError(t, os.MkdirAll(dir, 0755))

	ioutil.WriteFile(filepath.Join(dir, hookName), content, 0755)
	return func() {
		os.RemoveAll(dir)
	}
}
