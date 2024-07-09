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
	err := os.WriteFile("test_file_1.txt", []byte("hello world\n"), 0644)
	if err != nil {
		return err
	}

	err = os.MkdirAll("test_dir_1", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile("test_dir_1/test_file_2.txt", []byte("hello world test_dir_1/test_file_2.txt\n"), 0644)
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
		require.Equal(t, res, test_utils.FormatTreeChildren([]string{
			"tree 7073a74d71d9b2018918475aa6077630182b8acf",
			"author Author Name <author@example.com> 1720526831 -0300",
			"committer Author Name <author@example.com> 1720526831 -0300",
			"",
			"my commit",
		}))
	})

	t.Run("create correctly commit-tree with parent", func(t *testing.T) {
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

		// mock some second to first commit
		mockedSeconds := int64(1720529056)
		objects.Now = time.Unix(mockedSeconds, 0).In(loc)
		parentHash, _ := commands.CommitTree([]string{treeHash, "-m", "my commit"})

		// mock some second to second commit
		mockedSeconds = int64(1720526831)
		objects.Now = time.Unix(mockedSeconds, 0).In(loc)
		commitHash, err := commands.CommitTree([]string{treeHash, "-p", parentHash, "-m", "my commit with parent"})
		require.NoError(t, err)

		res, err := commands.CatFile([]string{"-p", commitHash})
		require.NoError(t, err)
		require.Equal(t, res, test_utils.FormatTreeChildren([]string{
			"tree 7073a74d71d9b2018918475aa6077630182b8acf",
			"parent 7c3a066dd451d494a6c03218b666697c689430e2",
			"author Author Name <author@example.com> 1720526831 -0300",
			"committer Author Name <author@example.com> 1720526831 -0300",
			"",
			"my commit with parent",
		}))
	})

	t.Run("error when parent hash is not a commit object", func(t *testing.T) {
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

		commitHash, err := commands.CommitTree([]string{treeHash, "-p", treeHash, "-m", "my commit"})
		require.Equal(t, "", commitHash)
		require.Equal(t, "fatal: "+treeHash+" is not a valid 'commit' object", err.Error())
	})

	t.Run("error when tree hash is not a tree object", func(t *testing.T) {
		dir, err := test_utils.GitInitSetup(t)
		defer func() {
			os.Chdir("../..")
			os.RemoveAll(dir)
		}()
		require.NoError(t, err)

		f := "my_blob_object.txt"
		err = os.WriteFile(f, []byte("I am not a tree object"), 0644)
		require.NoError(t, err)

		blobHash, err := commands.HashObject([]string{"-w", f})
		require.NoError(t, err)

		commitHash, err := commands.CommitTree([]string{blobHash, "-m", "my commit"})
		require.Equal(t, "", commitHash)
		require.Equal(t, "fatal: "+blobHash+" is not a valid 'tree' object", err.Error())
	})
}

func TestInvalidCommitTreeArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
		setup    func()
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
		{
			name:     "not valid tree object",
			args:     []string{"1234"},
			expected: "fatal: not a valid object name 1234",
		},
		{
			name: "not valid parent object",
			args: []string{"valid_hash", "-p", "12345"},
			setup: func() {
				commands.SHOULD_VALIDATE_TREE_HASH = false
			},
			expected: "fatal: not a valid object name 12345",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			res, err := commands.CommitTree(tc.args)
			require.NotNil(t, err)
			require.Equal(t, tc.expected, err.Error())
			require.Equal(t, "", res)
		})
	}
}
