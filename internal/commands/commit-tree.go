package commands

import (
	"fmt"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func CommitTree(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("No tree hash provided")
	}

	argsLen := len(args)
	var parentHash string
	var message string

	// TODO: improve this flags parsing
	if len(args) > 1 {
		firstFlag := args[1]
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
	}

	treeHash := args[0]

	if !objects.RepoCheckObjectId(treeHash) {
		return "", fmt.Errorf("fatal: not a valid object name %s", treeHash)
	}

	commitObject := objects.NewCommitObject(treeHash, parentHash, message)

	err := objects.SaveObject(commitObject.CommonObject)
	if err != nil {
		return "", fmt.Errorf("Error saving object -> %s", err)
	}

	return commitObject.Hash, nil
}
