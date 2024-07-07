package objects

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

type ObjectType string

const (
	Blob   ObjectType = "blob"
	Tree   ObjectType = "tree"
	Commit ObjectType = "commit"
	Tag    ObjectType = "tag"
)

type CommonObject struct {
	Type    ObjectType `json:"type"`
	Size    int        `json:"size"`
	Content string     `json:"content"`
}

func NewCommonObjectFromHash(hash string) CommonObject {
	content, err := readObjectContent(hash)
	if err != nil {
		panic(fmt.Errorf("Error reading object -> %s\n", err))
	}

	data, err := uncompressObjectContent(content)
	if err != nil {
		panic(fmt.Errorf("Error decompressing object -> %s\n", err))
	}

	return ParseCommonObject(data)
}

func ParseCommonObject(blob string) CommonObject {
	c := CommonObject{}

	typeIdx := strings.IndexByte(blob, ' ')
	c.Type = ObjectType(blob[:typeIdx])

	blob = blob[typeIdx+1:]
	idx := strings.IndexByte(blob, '\x00')

	size, err := strconv.Atoi(blob[:idx])
	if err != nil {
		panic(err)
	}
	c.Size = size

	c.Content = blob[idx+1:]

	return c
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
