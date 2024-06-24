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

func HashObject(args []string) {
	if len(args) < 4 {
		fmt.Fprintf(os.Stderr, "usage: mygit cat-file <object>\n")
		os.Exit(1)
	}

	flag := args[2]
	file := args[3]

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file -> %s\n", err)
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Error creating directory -> %s\n", err)
			os.Exit(1)
		}

		compressedBlob := bytes.Buffer{}
		w := zlib.NewWriter(&compressedBlob)
		w.Write([]byte(blob))
		w.Close()

		err = os.WriteFile(path.Join(objectPath, objectFile), compressedBlob.Bytes(), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file -> %s\n", err)
			os.Exit(1)
		}
		fmt.Println(sha1_hash)
	default:
		fmt.Fprintf(os.Stderr, "Unknown flag %s\n", flag)
		os.Exit(1)
	}
}
