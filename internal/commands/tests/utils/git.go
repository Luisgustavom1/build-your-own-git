package test_utils

import (
	"os"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
)

func GitInitSetup(t *testing.T) (dir string, err error) {
	err = os.MkdirAll("tests/tmp", 0755)
	if err != nil {
		return "", err
	}

	dir, err = os.MkdirTemp("tests/tmp", "test")
	if err != nil {
		return "", err
	}
	defer func() {
		os.Chdir("..")
		os.RemoveAll(dir)
	}()

	err = os.Chdir(dir)
	if err != nil {
		return "", err
	}

	_, err = commands.Init([]string{"mygit", "init"})
	return dir, err
}
