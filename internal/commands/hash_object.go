package commands

import (
	"fmt"
	"os"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func HashObject(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: mygit hash-object <object>\n")
	}

	flag := args[0]
	file := args[1]

	switch flag {
	case "-w":
		data, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("Error reading file -> %s\n", err)
		}

		blobObject := objects.NewBlobObject(data)
		err = objects.SaveObject(blobObject.CommonObject)
		if err != nil {
			return "", fmt.Errorf("Error saving object -> %s\n", err)
		}

		return blobObject.Hash, nil
	default:
		return "", fmt.Errorf("Unknown flag %s\n", flag)
	}
}
