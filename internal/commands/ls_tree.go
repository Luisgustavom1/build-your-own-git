package commands

import (
	"fmt"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func LsTree(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: mygit ls-tree <tree-ish>\n")
	}

	var objHash string
	var treeStringOpts objects.TreeStringOpts

	if args[0] == "--name-only" {
		treeStringOpts.NameOnly = true
		objHash = args[1]
	} else if len(args) == 2 {
		objHash = args[1]
	} else {
		objHash = args[0]
	}

	common := objects.NewCommonObjectFromHash(objHash)

	if common.Type != objects.Tree {
		return "", fmt.Errorf("not a tree object\n")
	}

	tree := objects.ParseTreeObject(common)

	return tree.String(treeStringOpts), nil
}
