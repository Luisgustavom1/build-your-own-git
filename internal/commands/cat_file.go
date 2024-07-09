package commands

import (
	"fmt"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func CatFile(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: mygit cat-file <object>\n")
	}

	flag := args[0]
	objHash := args[1]

	common := objects.NewObjectFromHash(objHash)

	switch flag {
	case "-t":
		return fmt.Sprintln(common.Type), nil
	case "-s":
		return fmt.Sprintln(common.Size), nil
	case "-p":
		if common.Type == objects.Tree {
			tree := objects.NewTreeObject(common)
			return tree.String(objects.TreeStringOpts{}), nil
		}
		if common.Type == objects.Commit {
			commit := objects.NewCommitObject(common)
			return commit.String(), nil
		}
		blob := objects.NewBlobObject(common)
		return blob.String(), nil
	default:
		return "", fmt.Errorf("Unknown flag %s\n", flag)
	}
}
