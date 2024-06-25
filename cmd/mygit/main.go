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

	commands.Orchestrator(os.Args)
}
