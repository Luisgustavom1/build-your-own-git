package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	test_utils "github.com/Luisgustavom1/build-your-own-git/internal/commands/tests/utils"
	"github.com/stretchr/testify/require"
)

func formatTreeChildren(strs []string) string {
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
			hashObj:  "d186cf338dd6da240c5c60a9f911dcd8e235b5c5",
			blobPath: "simple-tree-object",
			expectedLines: []string{
				"100644 blob b5b8fe9ad0f62425a834e50abf89b26f0a630902    catFile.go",
				"100644 blob 60724c6ef7823c90f20ed816dbaaeafe13915a44    hashObject.go",
				"100644 blob cd6154283fe7e083ba7baee6c4d06b786a4d36c4    init.go",
				"100644 blob 705b3f7027d5577c396bc5ec6fe0acdac5f83229    init_test.go",
				"100644 blob 23e9616617ff89d00f7599182d9b66d245f40ce1    orchestrator.go",
			},
		},
		{
			name:     "ls-tree with directory objects mode",
			hashObj:  "e6c55e07165517dad132ad455f6a8093d11512f0",
			blobPath: "simple-tree-object-with-dirs",
			expectedLines: []string{
				"040000 tree 99eb3cf514b0ae57271b15a86302b6e25b7f9493    test_dir_1",
				"040000 tree 0ef13a30a953dabfb813a0e9b24bb17dd6fa3ed9    test_dir_2",
				"100644 blob 95d09f2b10159347eece71399a7e2e907ea3df4f    test_file_1.txt",
			},
		},
		{
			name:     "ls-tree with --name-only flag",
			flag:     "--name-only",
			hashObj:  "d186cf338dd6da240c5c60a9f911dcd8e235b5c5",
			blobPath: "simple-tree-object",
			expectedLines: []string{
				"catFile.go",
				"hashObject.go",
				"init.go",
				"init_test.go",
				"orchestrator.go",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := test_utils.GitInitSetup(t)
			defer func() {
				os.Chdir("../..")
				os.RemoveAll(dir)
			}()
			require.NoError(t, err)

			objectPath := path.Join(".git/objects", tc.hashObj[:2])
			err = os.MkdirAll(objectPath, 0755)
			require.NoError(t, err)

			blob, err := os.ReadFile(path.Join("../../tests/fixtures", tc.blobPath))
			require.NoError(t, err)
			err = os.WriteFile(path.Join(objectPath, tc.hashObj[2:]), blob, 0644)
			require.NoError(t, err)

			args := []string{tc.flag, tc.hashObj}
			res, err := commands.LsTree(args)
			require.NoError(t, err)
			require.Equal(t, formatTreeChildren(tc.expectedLines), res)
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
