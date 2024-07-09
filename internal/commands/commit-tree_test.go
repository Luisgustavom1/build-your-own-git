package commands_test

import (
	"os"
	"testing"
	"time"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	test_utils "github.com/Luisgustavom1/build-your-own-git/internal/commands/tests/utils"
	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
	"github.com/stretchr/testify/require"
)

func setupCommitTreeFiles() error {
	err := os.WriteFile("test_file_1.txt", []byte("hello world"), 0644)
	if err != nil {
		return err
	}

	err = os.MkdirAll("test_dir_1", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile("test_dir_1/test_file_2.txt", []byte("hello world test_dir_1/test_file_2.txt"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func TestCommitTree(t *testing.T) {
	t.Run("create correctly commit-tree", func(t *testing.T) {
		dir, err := test_utils.GitInitSetup(t)
		defer func() {
			os.Chdir("../..")
			os.RemoveAll(dir)
		}()
		require.NoError(t, err)

		err = setupCommitTreeFiles()
		require.NoError(t, err)

		treeHash, err := commands.WriteTree([]string{})
		require.NoError(t, err)

		loc := time.FixedZone("-0300", -10800)
		mockedSeconds := int64(1720526831)
		objects.Now = time.Unix(mockedSeconds, 0).In(loc)

		commitHash, err := commands.CommitTree([]string{treeHash, "-m", "my commit"})
		require.NoError(t, err)
		require.Equal(t, "1dbca800d8f275e8fdf3e86470f7b76d5b09c3d4", commitHash)

		res, err := commands.CatFile([]string{"-p", commitHash})
		require.NoError(t, err)
		require.Contains(t, res, `tree 7073a74d71d9b2018918475aa6077630182b8acf
		author Author Name <author@example.com> 1720526831 -0300
		committer Author Name <author@example.com> 1720526831 -0300
		
		my commit`)
	})
}

func TestCommitTreeErrors(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "not a valid object name",
			args:     []string{"1234"},
			expected: "fatal: not a valid object name 1234",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := commands.CommitTree(tc.args)
			require.NotNil(t, err)
			require.Equal(t, tc.expected, err.Error())
			require.Equal(t, "", res)
		})
	}
}

func TestInvalidCommitTreeArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "commit tree with no flags",
			args:     []string{},
			expected: "No tree hash provided",
		},
		{
			name:     "commit tree without parent hash",
			args:     []string{"tree_hash", "-p"},
			expected: "No parent hash provided",
		},
		{
			name:     "commit tree without message",
			args:     []string{"tree_hash", "-m"},
			expected: "No message provided",
		},
		{
			name:     "commit tree without message",
			args:     []string{"tree_hash", "-p", "parent_hash", "-m"},
			expected: "No message provided",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := commands.CommitTree(tc.args)
			require.NotNil(t, err)
			require.Equal(t, tc.expected, err.Error())
			require.Equal(t, "", res)
		})
	}
}
