package commands

import (
	"fmt"
)

func LsTree(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("usage: mygit ls-tree <tree-ish>\n")
	}

	object := args[0]

	content, err := readObjectContent(object)
	if err != nil {
		return "", fmt.Errorf("Error reading object -> %s\n", err)
	}

	data, err := uncompressObjectContent(content)
	if err != nil {
		return "", fmt.Errorf("Error decompressing object -> %s\n", err)
	}

	common := parseCommonObject(data)

	if common.ttype != Tree {
		return "", fmt.Errorf("not a tree object\n")
	}

	tree := parseTreeObject(common)

	return tree.String(), nil
}
