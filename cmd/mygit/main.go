package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codecrafters-io/git-starter-go/cmd/mygit/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	e, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	executableDir := path.Dir(e)

	switch command := os.Args[1]; command {
	case "init":
		commands.Init(executableDir)
	case "cat-file":
		commands.CatFile(os.Args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
