package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codecrafters-io/git-starter-go/internal/commands"
)

func TestInit(t *testing.T) {
	t.Run("valid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp(".", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{"mygit", "init", dir}

		commands.Init(args)

		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Errorf("dir %s should be created", gitDir)
		}

		objectsDir := filepath.Join(gitDir, "objects")
		if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
			t.Errorf("dir %s should be created", objectsDir)
		}

		refsDir := filepath.Join(gitDir, "refs")
		if _, err := os.Stat(refsDir); os.IsNotExist(err) {
			t.Errorf("dir %s should be created", refsDir)
		}

		headFilePath := filepath.Join(gitDir, "HEAD")
		if _, err := os.Stat(headFilePath); os.IsNotExist(err) {
			t.Errorf("file %s should be created", headFilePath)
		}
	})

	t.Run("invalid arguments", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		args := []string{"mygit", "init"}

		commands.Init(args)

		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
			t.Errorf("Directory %s was created", gitDir)
		}

		headFilePath := filepath.Join(gitDir, "HEAD")
		if _, err := os.Stat(headFilePath); !os.IsNotExist(err) {
			t.Errorf("File %s was created", headFilePath)
		}
	})
}
