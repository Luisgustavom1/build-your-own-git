package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	test_utils "github.com/Luisgustavom1/build-your-own-git/internal/commands/tests/utils"
	"github.com/stretchr/testify/require"
)

func TestCatFile(t *testing.T) {
	testCases := []struct {
		name              string
		flag              string
		object            string
		objectContentPath string
		expected          string
	}{
		{
			name:              "object of blob type with flag -t",
			flag:              "-t",
			object:            "3b18e512dba79e4c8300dd08aeb37f8e728b8dad",
			objectContentPath: "hash-object-hello-world",
			expected:          "blob\n",
		},
		{
			name:              "object of blob type with flag -s",
			flag:              "-s",
			object:            "3b18e512dba79e4c8300dd08aeb37f8e728b8dad",
			objectContentPath: "hash-object-hello-world",
			expected:          "12\n",
		},
		{
			name:              "object of blob type with flag -p",
			flag:              "-p",
			object:            "3b18e512dba79e4c8300dd08aeb37f8e728b8dad",
			objectContentPath: "hash-object-hello-world",
			expected:          "hello world\n",
		},
		{
			name:              "object of tree type with flag -t",
			flag:              "-t",
			object:            "d186cf338dd6da240c5c60a9f911dcd8e235b5c5",
			objectContentPath: "simple-tree-object",
			expected:          "tree\n",
		},
		{
			name:              "object of tree type",
			flag:              "-p",
			object:            "d186cf338dd6da240c5c60a9f911dcd8e235b5c5",
			objectContentPath: "simple-tree-object",
			expected: test_utils.FormatTreeChildren([]string{
				"100644 blob b5b8fe9ad0f62425a834e50abf89b26f0a630902    catFile.go",
				"100644 blob 60724c6ef7823c90f20ed816dbaaeafe13915a44    hashObject.go",
				"100644 blob cd6154283fe7e083ba7baee6c4d06b786a4d36c4    init.go",
				"100644 blob 705b3f7027d5577c396bc5ec6fe0acdac5f83229    init_test.go",
				"100644 blob 23e9616617ff89d00f7599182d9b66d245f40ce1    orchestrator.go",
			}),
		},
		{
			name:              "object of commit type",
			flag:              "-p",
			object:            "1dbca800d8f275e8fdf3e86470f7b76d5b09c3d4",
			objectContentPath: "simple-commit-object",
			expected: test_utils.FormatTreeChildren([]string{
				"tree 7073a74d71d9b2018918475aa6077630182b8acf",
				"author Author Name <author@example.com> 1720526831 -0300",
				"committer Author Name <author@example.com> 1720526831 -0300",
				"",
				"my commit",
			}),
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

			objectPath := path.Join(".git/objects", tc.object[:2])
			err = os.MkdirAll(objectPath, 0755)
			require.NoError(t, err)

			blob, err := os.ReadFile(path.Join("../../tests/fixtures", tc.objectContentPath))
			require.NoError(t, err)
			err = os.WriteFile(path.Join(objectPath, tc.object[2:]), blob, 0644)
			require.NoError(t, err)

			args := []string{tc.flag, tc.object}
			res, err := commands.CatFile(args)
			require.NoError(t, err)
			require.Equal(t, tc.expected, res)
		})
	}

	t.Run("invalid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{"-p"}

		_, err = commands.CatFile(args)
		require.EqualErrorf(t, err, "usage: mygit cat-file <object>\n", "Invalid args")
	})
}
