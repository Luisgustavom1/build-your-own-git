package commands

import "fmt"

func CommitTree(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("No tree hash provided")
	}

	firstFlag := args[1]
	argsLen := len(args)
	var parentHash string
	var message string

	// TODO: improve this flags parsing
	if firstFlag == "-p" {
		if argsLen < 3 {
			return "", fmt.Errorf("No parent hash provided")
		}
		parentHash = args[2]
	}

	if firstFlag == "-m" {
		if argsLen < 3 {
			return "", fmt.Errorf("No message provided")
		}
		message = args[2]
	}

	if argsLen > 3 && args[3] == "-m" {
		if argsLen < 5 {
			return "", fmt.Errorf("No message provided")
		}
		message = args[4]
	}

	treeHash := args[0]

	fmt.Println("treeHash: ", treeHash)
	fmt.Println("parentHash: ", parentHash)
	fmt.Println("message: ", message)

	return "", nil
}
