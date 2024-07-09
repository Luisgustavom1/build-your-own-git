package commands

import (
	"fmt"
	"os"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func WriteTree(args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("usage: mygit write-tree\n")
	}

	dirs, err := os.ReadDir("./")
	if err != nil {
		return "", fmt.Errorf("Error reading directory -> %s\n", err)
	}

	tree := objects.TreeObject{}
	err = objects.GenerateTreeObjectFromDirs(dirs, ".", &tree)
	if err != nil {
		return "", fmt.Errorf("Error generating tree object -> %s\n", err)
	}

	err = objects.SaveNodes(tree)
	if err != nil {
		return "", err
	}

	return tree.Value.Hash, nil
}
