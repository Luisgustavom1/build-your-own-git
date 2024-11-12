package commands

import (
	"fmt"
	"os"
)

// TODO: use some cli framework
func Orchestrator(args []string) {
	var res string
	var err error
	command := args[1]
	args = args[2:]
	switch command {
	case "init":
		res, err = Init(args)
	case "cat-file":
		res, err = CatFile(args)
	case "hash-object":
		res, err = HashObject(args)
	case "ls-tree":
		res, err = LsTree(args)
	case "commit-tree":
		res, err = CommitTree(args)
	case "clone":
		res, err = Clone(args)
	default:
		err = fmt.Errorf("Unknown command %s\n", command)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(res)
}
