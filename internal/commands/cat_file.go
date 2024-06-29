package commands

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type ObjectType string

const (
	Blob   ObjectType = "blob"
	Tree   ObjectType = "tree"
	Commit ObjectType = "commit"
	Tag    ObjectType = "tag"
)

type ObjectModes string

const (
	RegularFileMode    = "100644"
	ExecutableFileMode = "100755"
	SymbolicLinkMode   = "120000"
	DirectoryMode      = "040000"
)

type CommonObject struct {
	ttype   ObjectType
	size    string
	content string
}

type BlobObject struct {
	CommonObject
	data string
}

type TreeObjectContent struct {
	mode  string
	ttype ObjectType
	name  string
	hash  string
}

type TreeObject struct {
	CommonObject
	contents []TreeObjectContent
}

type Object interface {
	String() string
}

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

	common := parseCommonObject(data)

	switch flag {
	case "-t":
		return fmt.Sprintln(common.ttype), nil
	case "-s":
		return fmt.Sprintln(common.size), nil
	case "-p":
		content := parseObjectContent(common)
		return content.String(), nil
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

func parseCommonObject(blob string) CommonObject {
	c := CommonObject{}

	typeIdx := strings.Index(blob, " ")
	c.ttype = ObjectType(blob[:typeIdx])

	blob = blob[typeIdx+1:]
	idx := strings.Index(blob, "\x00")
	c.size = blob[:idx]

	c.content = blob[idx+1:]

	return c
}

func parseObjectContent(object CommonObject) Object {
	switch object.ttype {
	case Blob:
		return parseBlobObject(object)
	case Tree:
		return parseTreeObject(object)
	default:
		return nil
	}
}

func (b BlobObject) String() string {
	return b.data
}

func parseBlobObject(object CommonObject) Object {
	blob := BlobObject{CommonObject: object}
	blob.data = object.content
	return blob
}

func (t TreeObject) String() string {
	b := strings.Builder{}
	for _, content := range t.contents {
		b.WriteString(fmt.Sprintf("%s %s %s    %s\n", content.mode, content.ttype, content.hash, content.name))
	}
	return b.String()
}

func parseTreeObject(object CommonObject) Object {
	nodes := make([]TreeObjectContent, 0)
	last := TreeObjectContent{}

	idx := 0
	s := strings.Builder{}

	for idx < len(object.content) {
		// <mode> <name>\0<20_byte_sha>
		// <mode> <name>\0<20_byte_sha>
		str := object.content[idx]
		if string(str) == " " {
			last.mode = s.String()
			s.Reset()

			// we cover only this mode for now
			if last.mode == DirectoryMode {
				last.ttype = Tree
			} else {
				last.ttype = Blob
			}

			nullIdx := strings.Index(object.content[idx:], "\x00")

			name := object.content[idx+1 : idx+nullIdx]
			last.name = name

			idx += nullIdx + 1

			reader := strings.NewReader(object.content[idx:])

			hash := make([]byte, 20)
			_, _ = io.ReadAtLeast(reader, hash, 20)

			last.hash = hex.EncodeToString(hash[:])

			// 20 bytes (hash) + 1 byte (jump to next node)
			idx += 20

			nodes = append(nodes, last)
			last = TreeObjectContent{}
			continue
		}

		_, err := s.WriteString(string(str))
		if err != nil {
			panic(err)
		}
		idx++
	}

	return TreeObject{
		CommonObject: object,
		contents:     nodes,
	}
}

func parseTreeObjectContent(nodeStr string) TreeObjectContent {
	node := TreeObjectContent{}

	idx := strings.Index(nodeStr, " ")
	node.mode = nodeStr[:idx]

	nodeStr = nodeStr[idx+1:]
	idx = strings.Index(nodeStr, " ")
	node.ttype = ObjectType(nodeStr[:idx])

	nodeStr = nodeStr[idx+1:]
	idx = strings.Index(nodeStr, "\x00")
	node.hash = nodeStr[:idx]

	return node
}
