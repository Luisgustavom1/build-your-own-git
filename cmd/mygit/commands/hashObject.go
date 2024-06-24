package commands

import (
	"fmt"
	"os"
)

func HashObject(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: mygit cat-file <object>\n")
		os.Exit(1)
	}

	flag := args[2]
	object := args[3]

	if object == "" {
		fmt.Fprintf(os.Stderr, "usage: mygit cat-file <object>\n")
		return
	}

	switch flag {
	case "-w":
		fmt.Println("write a file")
	default:
		fmt.Fprintf(os.Stderr, "Unknown flag %s\n", flag)
		os.Exit(1)
	}
}
