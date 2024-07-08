package commands

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/Luisgustavom1/build-your-own-git/internal/objects"
)

const GIT_FOLDER = ".git"

func WriteTree(args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("usage: mygit write-tree\n")
	}

	dirs, err := os.ReadDir("./")
	if err != nil {
		return "", fmt.Errorf("Error reading directory -> %s\n", err)
	}

	tree := objects.TreeObject{}
	err = generateTreeObject(dirs, ".", &tree)
	if err != nil {
		return "", fmt.Errorf("Error generating tree object -> %s\n", err)
	}

	err = SaveNodes(tree)
	if err != nil {
		return "", err
	}

	return tree.Value.Hash, nil
}

// TODO: move to tree package
func generateTreeObject(dirs []fs.DirEntry, wd string, tree *objects.TreeObject) error {
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

			node := &objects.TreeObject{}
			err = generateTreeObject(subDir, fullPath, node)
			if err != nil {
				return fmt.Errorf("Error generating tree object -> %s\n", err)
			}

			tree.Children = append(tree.Children, node)
			continue
		}

		node := &objects.TreeObject{}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("Error reading file -> %s\n", err)
		}

		blobObject := objects.NewBlobObject(data)

		node.Value.Mode = objects.RegularFileMode
		node.Value.Name = dir.Name()
		node.Value.CommonObject = blobObject.CommonObject

		tree.Children = append(tree.Children, node)
	}

	s, c, h := createTreeSha1Hash(tree)
	tree.Value.Mode = objects.DirectoryMode
	tree.Value.Type = objects.Tree
	tree.Value.Size = s
	tree.Value.Content = c
	tree.Value.Hash = h
	tree.Value.Name = wd

	return nil
}

func createTreeSha1Hash(tree *objects.TreeObject) (size int, content string, hash string) {
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

func SaveNodes(tree objects.TreeObject) error {
	var err error
	err = objects.SaveObject(tree.Value.CommonObject)
	if err != nil {
		return fmt.Errorf("Error saving object %+v -> %s\n", tree.Value, err)
	}
	for _, child := range tree.Children {
		err = objects.SaveObject(child.Value.CommonObject)
		if err != nil {
			return fmt.Errorf("Error saving object %+v -> %s\n", child.Value, err)
		}
	}
	return nil
}
