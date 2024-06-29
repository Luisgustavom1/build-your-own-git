package commands

import (
	"fmt"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func LsTree(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: mygit ls-tree <tree-ish>\n")
	}

	var object string
	var treeStringOpts objects.TreeStringOpts

	if args[0] == "--name-only" {
		treeStringOpts.NameOnly = true
		object = args[1]
	} else if len(args) == 2 {
		object = args[1]
	} else {
		object = args[0]
	}

	content, err := readObjectContent(object)
	if err != nil {
		return "", fmt.Errorf("Error reading object -> %s\n", err)
	}

	data, err := uncompressObjectContent(content)
	if err != nil {
		return "", fmt.Errorf("Error decompressing object -> %s\n", err)
	}

	common := objects.ParseCommonObject(data)

	if common.Type != objects.Tree {
		return "", fmt.Errorf("not a tree object\n")
	}

	return parseAndPrintTree(common, treeStringOpts), nil
}

func parseAndPrintTree(common objects.CommonObject, opts objects.TreeStringOpts) string {
	tree := objects.ParseTreeObject(common)
	return objects.TreeObjToString(tree, opts)
}
