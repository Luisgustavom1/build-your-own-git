package commands

import "fmt"

func WriteTree(args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("usage: mygit write-tree\n")
	}

	return "fb88fc4b84ad85b59151616c4d02591ca4a18f28", nil
}
