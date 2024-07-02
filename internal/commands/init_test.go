package commands_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("valid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp("./tmp", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{dir}

		res, err := commands.Init(args)
		require.NoError(t, err)

		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Errorf("dir %s should be created", gitDir)
		}

		objectsDir := filepath.Join(gitDir, "objects")
		_, err = os.Stat(objectsDir)
		require.False(t, os.IsNotExist(err), fmt.Sprintf("dir %s should be created", objectsDir))

		refsDir := filepath.Join(gitDir, "refs")
		_, err = os.Stat(refsDir)
		require.False(t, os.IsNotExist(err), fmt.Sprintf("dir %s should be created", refsDir))

		headFilePath := filepath.Join(gitDir, "HEAD")
		_, err = os.Stat(headFilePath)
		require.False(t, os.IsNotExist(err), fmt.Sprintf("file %s should be created", headFilePath))

		require.Contains(t, res, "Initialized empty Git repository in")
	})
}
