package commands

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

func HashObject(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: mygit hash-object <object>\n")
	}

	flag := args[0]
	file := args[1]

	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file -> %s\n", err)
	}

	switch flag {
	case "-w":
		blob := fmt.Sprintf("blob %d\000%s", len(data), data)
		h := sha1.Sum([]byte(blob))
		sha1_hash := hex.EncodeToString(h[:])
		objectPath := path.Join(".git/objects", string(sha1_hash[:2]))
		objectFile := sha1_hash[2:]

		err := os.MkdirAll(objectPath, 0755)
		if err != nil {
			return "", fmt.Errorf("Error creating directory -> %s\n", err)
		}

		compressedBlob := bytes.Buffer{}
		w := zlib.NewWriter(&compressedBlob)
		w.Write([]byte(blob))
		w.Close()

		err = os.WriteFile(path.Join(objectPath, objectFile), compressedBlob.Bytes(), 0644)
		if err != nil {
			return "", fmt.Errorf("Error writing file -> %s\n", err)
		}

		return fmt.Sprintln(sha1_hash), nil
	default:
		return "", fmt.Errorf("Unknown flag %s\n", flag)
	}
}
