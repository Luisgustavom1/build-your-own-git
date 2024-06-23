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

func CatFile(args []string) {
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

	content, err := readObjectContent(object)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object -> %s\n", err)
		os.Exit(1)
	}

	blob := uncompressObjectContent(content)

	objBlob := parseObjectBlob(blob)

	switch flag {
	case "-t":
		fmt.Println(objBlob.ttype)
	case "-s":
		fmt.Println(objBlob.size)
	case "-p":
		fmt.Print(objBlob.data)
	default:
		fmt.Fprintf(os.Stderr, "Unknown flag %s\n", flag)
		os.Exit(1)
	}
}

func readObjectContent(object string) ([]byte, error) {
	return os.ReadFile(path.Join(".git/objects", object[:2], object[2:]))
}

func uncompressObjectContent(content []byte) string {
	buff := bytes.NewBuffer([]byte(content))
	r, err := zlib.NewReader(buff)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in decompressing object -> %s\n", err)
		os.Exit(1)
	}
	defer r.Close()

	blob := &bytes.Buffer{}
	io.Copy(blob, r)
	return blob.String()
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
