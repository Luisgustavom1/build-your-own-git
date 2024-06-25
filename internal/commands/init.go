package commands

import (
	"fmt"
	"os"
	"path"
)

func Init(args []string) {
	var parentDir string
	if len(args) == 3 {
		parentDir = args[2]
	}

	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(path.Join(parentDir, dir), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(path.Join(parentDir, ".git/HEAD"), headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %s\n", err)
	}

	fmt.Println("Initialized empty Git repository in", path.Join(wd, parentDir, ".git"))
}
