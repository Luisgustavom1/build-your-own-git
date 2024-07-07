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

	common := objects.NewCommonObjectFromHash(objHash)

	switch flag {
	case "-t":
		return fmt.Sprintln(common.Type), nil
	case "-s":
		return fmt.Sprintln(common.Size), nil
	case "-p":
		if common.Type == objects.Tree {
			tree := objects.ParseTreeObject(common)
			return tree.String(objects.TreeStringOpts{}), nil
		}
		blob := objects.ParseBlobObject(common)
		return blob.String(), nil
	default:
		return "", fmt.Errorf("Unknown flag %s\n", flag)
	}
}
