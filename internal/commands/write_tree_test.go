package commands_test

import (
	"os"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	test_utils "github.com/Luisgustavom1/build-your-own-git/internal/commands/tests/utils"
	"github.com/stretchr/testify/require"
)

func setupTreeFiles() error {
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

	err = os.MkdirAll("test_dir_2", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile("test_dir_2/test_file_3.txt", []byte("hello world test_dir_2/test_file_3.txt"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func TestWriteTree(t *testing.T) {
	t.Run("create correctly write-tree", func(t *testing.T) {
		dir, err := test_utils.GitInitSetup(t)
		defer func() {
			os.Chdir("../..")
			os.RemoveAll(dir)
		}()
		require.NoError(t, err)

		err = setupTreeFiles()
		require.NoError(t, err)

		res, err := commands.WriteTree([]string{})
		require.NoError(t, err)

		expectedHash := "e6c55e07165517dad132ad455f6a8093d11512f0"
		require.Equal(t, expectedHash, res)

		res, err = commands.LsTree([]string{expectedHash})
		require.NoError(t, err)

		tree := test_utils.FormatTreeChildren([]string{
			"040000 tree 99eb3cf514b0ae57271b15a86302b6e25b7f9493    test_dir_1",
			"040000 tree 0ef13a30a953dabfb813a0e9b24bb17dd6fa3ed9    test_dir_2",
			"100644 blob 95d09f2b10159347eece71399a7e2e907ea3df4f    test_file_1.txt",
		})
		require.Equal(t, tree, res)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{"invalid-arg-1", "invalid-arg-2"}

		_, err = commands.WriteTree(args)
		require.EqualErrorf(t, err, "usage: mygit write-tree\n", "Invalid args")
	})
}
