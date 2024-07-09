package commands

import (
	"fmt"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

// TODO: move this to some env var
var SHOULD_VALIDATE_TREE_HASH = true

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

	if err := validateTreeObj(treeHash); err != nil {
		return "", err
	}

	if err := validateParentObj(parentHash); err != nil {
		return "", err
	}

	data := objects.GenerateCommitData(treeHash, parentHash, message)
	commitObject := objects.NewCommitFromData(data)
	err := commitObject.Write()
	if err != nil {
		return "", fmt.Errorf("Error saving object -> %s", err)
	}

	return commitObject.Hash, nil
}

func validateTreeObj(hash string) error {
	if !SHOULD_VALIDATE_TREE_HASH {
		return nil
	}

	if hash == "" {
		return fmt.Errorf("No tree hash provided")
	}

	if !objects.RepoCheckObjectId(hash) {
		return fmt.Errorf("fatal: not a valid object name %s", hash)
	}

	treeObject := objects.NewObjectFromHash(hash)
	if treeObject.Type != objects.Tree {
		return fmt.Errorf("fatal: %s is not a valid 'tree' object", hash)
	}

	return nil
}

func validateParentObj(hash string) error {
	if hash == "" {
		return nil
	}

	if !objects.RepoCheckObjectId(hash) {
		return fmt.Errorf("fatal: not a valid object name %s", hash)
	}

	parentObject := objects.NewObjectFromHash(hash)
	if parentObject.Type != objects.Commit {
		return fmt.Errorf("fatal: %s is not a valid 'commit' object", hash)
	}

	return nil
}
