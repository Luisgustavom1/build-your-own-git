package commands

import (
	"fmt"
	"os"
)

func Orchestrator(args []string) {
	var res string
	var err error

	switch command := args[1]; command {
	case "init":
		res, err = Init(args)
	case "cat-file":
		res, err = CatFile(args)
	case "hash-object":
		res, err = HashObject(args)
	default:
		err = fmt.Errorf("Unknown command %s\n", command)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(res)
}
