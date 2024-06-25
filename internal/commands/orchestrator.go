package commands

import (
	"fmt"
	"os"
)

func Orchestrator(args []string) {
	var err error

	switch command := args[1]; command {
	case "init":
		err = Init(args)
	case "cat-file":
		err = CatFile(args)
	case "hash-object":
		err = HashObject(args)
	default:
		err = fmt.Errorf("Unknown command %s\n", command)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
