package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		commands.Init(os.Args)
	case "cat-file":
		commands.CatFile(os.Args)
	case "hash-object":
		commands.HashObject(os.Args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
