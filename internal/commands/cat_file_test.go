package commands_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	test_utils "github.com/Luisgustavom1/build-your-own-git/internal/commands/tests/utils"
	"github.com/stretchr/testify/require"
)

func TestCatFile(t *testing.T) {
	testCases := []struct {
		name     string
		flag     string
		hashObj  string
		blobPath string
		expected string
	}{
		{
			name:     "with flag -t",
			flag:     "-t",
			hashObj:  "3b18e512dba79e4c8300dd08aeb37f8e728b8dad",
			blobPath: "hash-object-hello-world",
			expected: "blob",
		},
		{
			name:     "with flag -s",
			flag:     "-s",
			hashObj:  "3b18e512dba79e4c8300dd08aeb37f8e728b8dad",
			blobPath: "hash-object-hello-world",
			expected: "12",
		},
		{
			name:     "with flag -p",
			flag:     "-p",
			hashObj:  "3b18e512dba79e4c8300dd08aeb37f8e728b8dad",
			blobPath: "hash-object-hello-world",
			expected: "hello world",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := test_utils.GitInitSetup(t)
			require.NoError(t, err)

			objectPath := path.Join(".git/objects", tc.hashObj[:2])
			err = os.MkdirAll(objectPath, 0755)
			require.NoError(t, err)

			blob, err := os.ReadFile(path.Join("tests/fixtures", tc.blobPath))
			require.NoError(t, err)
			err = os.WriteFile(path.Join(objectPath, tc.hashObj[2:]), blob, 0644)
			require.NoError(t, err)

			args := []string{"mygit", "cat-file", tc.flag, tc.hashObj}
			res, err := commands.CatFile(args)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintln(tc.expected), res)
		})
	}

	t.Run("invalid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{"mygit", "cat-file"}

		_, err = commands.CatFile(args)
		require.EqualErrorf(t, err, "usage: mygit cat-file <object>\n", "Invalid args")
	})
}
