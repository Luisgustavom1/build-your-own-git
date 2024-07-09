package objects

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
)

type ObjectModes string

const (
	RegularFileMode    = "100644"
	ExecutableFileMode = "100755"
	SymbolicLinkMode   = "120000"
	DirectoryMode      = "40000"
)

type TreeObjectNode struct {
	CommonObject
	Mode string `json:"mode"`
	Name string `json:"name"`
}

type TreeObject struct {
	Value    TreeObjectNode `json:"value"`
	Children []*TreeObject  `json:"children"`
}

type TreeStringOpts struct {
	NameOnly bool
}

const GIT_FOLDER = ".git"

func ParseTreeObject(object CommonObject) TreeObject {
	nodes := make([]*TreeObject, 0)
	last := &TreeObject{}

	idx := 0
	s := strings.Builder{}

	for idx < len(object.Data) {
		// <mode> <name>\0<20_byte_sha>
		// <mode> <name>\0<20_byte_sha>
		str := object.Data[idx]
		if string(str) == " " {
			last.Value.Mode = s.String()
			s.Reset()

			// we cover only this mode for now
			if last.Value.Mode == DirectoryMode {
				last.Value.Type = Tree
			} else {
				last.Value.Type = Blob
			}

			nullIdx := strings.IndexByte(object.Data[idx:], '\x00')

			name := object.Data[idx+1 : idx+nullIdx]
			last.Value.Name = name

			idx += nullIdx + 1

			reader := strings.NewReader(object.Data[idx:])

			hash := make([]byte, 20)
			_, _ = io.ReadAtLeast(reader, hash, 20)

			last.Value.Hash = hex.EncodeToString(hash[:])

			// 20 bytes (hash)
			idx += 20

			nodes = append(nodes, last)
			last = &TreeObject{}
			continue
		}

		err := s.WriteByte(str)
		if err != nil {
			panic(err)
		}
		idx++
	}

	return TreeObject{
		Value: TreeObjectNode{
			CommonObject: object,
		},
		Children: nodes,
	}
}

// TODO: review and improve this code
func GenerateTreeObjectFromDirs(dirs []fs.DirEntry, wd string, tree *TreeObject) error {
	for _, dir := range dirs {
		if dir.Name() == GIT_FOLDER {
			continue
		}

		fullPath := path.Join(wd, dir.Name())

		if dir.IsDir() {
			subDir, err := os.ReadDir(dir.Name())
			if err != nil {
				return fmt.Errorf("Error reading directory -> %s\n", err)
			}

			node := &TreeObject{}
			err = GenerateTreeObjectFromDirs(subDir, fullPath, node)
			if err != nil {
				return fmt.Errorf("Error generating tree object -> %s\n", err)
			}

			tree.Children = append(tree.Children, node)
			continue
		}

		node := &TreeObject{}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("Error reading file -> %s\n", err)
		}

		blobObject := NewBlobObject(data)

		node.Value.Mode = RegularFileMode
		node.Value.Name = dir.Name()
		node.Value.CommonObject = blobObject.CommonObject

		tree.Children = append(tree.Children, node)
	}

	s, c, h := createTreeSha1Hash(tree)
	tree.Value.Mode = DirectoryMode
	tree.Value.Type = Tree
	tree.Value.Size = s
	tree.Value.Content = c
	tree.Value.Hash = h
	tree.Value.Name = wd

	return nil
}

func createTreeSha1Hash(tree *TreeObject) (size int, content string, hash string) {
	nodesContent := strings.Builder{}
	for _, child := range tree.Children {
		hexHash, err := hex.DecodeString(child.Value.Hash)
		if err != nil {
			panic(err)
		}
		nodesContent.WriteString(fmt.Sprintf("%s %s\x00%s", child.Value.Mode, child.Value.Name, hexHash))
	}
	nodesContentStr := nodesContent.String()
	size = len(nodesContentStr)

	c := fmt.Sprintf("tree %d\000%s", size, nodesContentStr)
	h := sha1.Sum([]byte(c))
	hash = hex.EncodeToString(h[:])

	return size, c, hash
}

func SaveNodes(tree TreeObject) error {
	var err error
	err = SaveObject(tree.Value.CommonObject)
	if err != nil {
		return fmt.Errorf("Error saving object %+v -> %s\n", tree.Value, err)
	}
	for _, child := range tree.Children {
		err = SaveObject(child.Value.CommonObject)
		if err != nil {
			return fmt.Errorf("Error saving object %+v -> %s\n", child.Value, err)
		}
	}
	return nil
}

func (t TreeObject) String(opts TreeStringOpts) string {
	b := strings.Builder{}
	for _, content := range t.Children {
		if opts.NameOnly {
			b.WriteString(fmt.Sprintf("%s\n", content.Value.Name))
		} else {
			b.WriteString(fmt.Sprintf("%06s %s %s    %s\n", content.Value.Mode, content.Value.Type, content.Value.Hash, content.Value.Name))
		}
	}
	return b.String()
}
