package commands

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type ObjectBlob struct {
	ttype string
	size  string
	data  string
}

func CatFile(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("usage: mygit cat-file <object>\n")
	}

	flag := args[2]
	object := args[3]

	if object == "" {
		return "", fmt.Errorf("usage: mygit cat-file <object>\n")
	}

	content, err := readObjectContent(object)
	if err != nil {
		return "", fmt.Errorf("Error reading object -> %s\n", err)
	}

	blob, err := uncompressObjectContent(content)
	if err != nil {
		return "", fmt.Errorf("Error decompressing object -> %s\n", err)
	}

	objBlob := parseObjectBlob(blob)

	switch flag {
	case "-t":
		return fmt.Sprintln(objBlob.ttype), nil
	case "-s":
		return fmt.Sprintln(objBlob.size), nil
	case "-p":
		return objBlob.data, nil
	default:
		return "", fmt.Errorf("Unknown flag %s\n", flag)
	}
}

func readObjectContent(object string) ([]byte, error) {
	return os.ReadFile(path.Join(".git/objects", object[:2], object[2:]))
}

func uncompressObjectContent(content []byte) (string, error) {
	buff := bytes.NewBuffer([]byte(content))
	r, err := zlib.NewReader(buff)
	if err != nil {
		return "", fmt.Errorf("Error in decompressing object -> %s\n", err)
	}
	defer r.Close()

	blob := &bytes.Buffer{}
	io.Copy(blob, r)
	return blob.String(), nil
}

func parseObjectBlob(blob string) ObjectBlob {
	objectBlob := ObjectBlob{}
	index := strings.Index(blob, " ")
	objectBlob.ttype = blob[:index]

	blob = blob[index+1:]
	index = strings.Index(blob, "\x00")
	objectBlob.size = blob[:index]

	blob = blob[index+1:]
	objectBlob.data = blob

	return objectBlob
}
