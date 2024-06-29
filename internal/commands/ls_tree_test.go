package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	test_utils "github.com/Luisgustavom1/build-your-own-git/internal/commands/tests/utils"
	"github.com/stretchr/testify/require"
)

func formatNodeContentsLine(strs []string) string {
	var res string
	for _, s := range strs {
		res += s + "\n"
	}
	return res
}

func TestLsTree(t *testing.T) {
	testCases := []struct {
		name          string
		flag          string
		hashObj       string
		blobPath      string
		expectedLines []string
	}{
		{
			name:     "ls-tree with valid tree",
			flag:     "-p",
			hashObj:  "d186cf338dd6da240c5c60a9f911dcd8e235b5c5",
			blobPath: "commands-dir-tree-object",
			expectedLines: []string{
				"100644 blob b5b8fe9ad0f62425a834e50abf89b26f0a630902    catFile.go",
				"100644 blob 60724c6ef7823c90f20ed816dbaaeafe13915a44    hashObject.go",
				"100644 blob cd6154283fe7e083ba7baee6c4d06b786a4d36c4    init.go",
				"100644 blob 705b3f7027d5577c396bc5ec6fe0acdac5f83229    init_test.go",
				"100644 blob 23e9616617ff89d00f7599182d9b66d245f40ce1    orchestrator.go",
			},
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

			args := []string{tc.flag, tc.hashObj}
			res, err := commands.CatFile(args)
			require.NoError(t, err)
			require.Equal(t, formatNodeContentsLine(tc.expectedLines), res)
		})
	}

	t.Run("invalid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{}

		_, err = commands.LsTree(args)
		require.EqualErrorf(t, err, "usage: mygit ls-tree <tree-ish>\n", "Invalid args")
	})
}
