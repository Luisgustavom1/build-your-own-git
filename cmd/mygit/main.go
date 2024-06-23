package main

import (
	"fmt"
	"os"
	"path"
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
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(path.Join(executableDir, dir), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized empty Git repository in", path.Join(executableDir, ".git"))
	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file <object>\n")
			os.Exit(1)
		}

		object := os.Args[2]
		fmt.Println(object)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
