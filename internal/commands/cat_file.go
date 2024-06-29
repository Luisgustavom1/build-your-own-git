package commands

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

func CatFile(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: mygit cat-file <object>\n")
	}

	flag := args[0]
	object := args[1]

	content, err := readObjectContent(object)
	if err != nil {
		return "", fmt.Errorf("Error reading object -> %s\n", err)
	}

	data, err := uncompressObjectContent(content)
	if err != nil {
		return "", fmt.Errorf("Error decompressing object -> %s\n", err)
	}

	common := objects.ParseCommonObject(data)

	switch flag {
	case "-t":
		return fmt.Sprintln(common.Type), nil
	case "-s":
		return fmt.Sprintln(common.Size), nil
	case "-p":
		if common.Type == objects.Tree {
			return parseAndPrintTree(common, objects.TreeStringOpts{}), nil
		}
		blob := objects.ParseBlobObject(common)
		return objects.BlobObjToString(blob), nil
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
