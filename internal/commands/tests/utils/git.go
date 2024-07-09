package test_utils

import (
	"os"
	"testing"

	"github.com/Luisgustavom1/build-your-own-git/internal/commands"
)

func GitInitSetup(t *testing.T) (dir string, err error) {
	err = os.MkdirAll("./tmp", 0755)
	if err != nil {
		return "", err
	}

	dir, err = os.MkdirTemp("./tmp", "test")
	if err != nil {
		return "", err
	}

	err = os.Chdir(dir)
	if err != nil {
		return "", err
	}

	_, err = commands.Init([]string{"mygit", "init"})
	return dir, err
}

func FormatTreeChildren(strs []string) string {
	var res string
	for _, s := range strs {
		res += s + "\n"
	}
	return res
}
