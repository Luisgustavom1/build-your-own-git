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

	hash := hex.EncodeToString([]byte(tree.Value.Hash))
	return hash, nil
}

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

			node := objects.TreeObject{}
			err = generateTreeObject(subDir, fullPath, &node)
			if err != nil {
				return fmt.Errorf("Error generating tree object -> %s\n", err)
			}

			tree.Children = append(tree.Children, node)
			continue
		}

		node := objects.TreeObject{}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("Error reading file -> %s\n", err)
		}

		blobObject := objects.NewBlobObject(data)

		node.Type = objects.Blob
		node.Size = len(data)
		node.Content = string(data)

		node.Value.Mode = objects.RegularFileMode
		node.Value.Name = dir.Name()

		hexHash, err := hex.DecodeString(blobObject.Hash)
		if err != nil {
			return fmt.Errorf("Error decoding hash -> %s\n", err)
		}
		node.Value.Hash = string(hexHash)

		tree.Children = append(tree.Children, node)
	}

	tree.Type = objects.Tree
	tree.Value.Name = wd
	tree.Value.Mode = objects.DirectoryMode
	createTreeSha1Hash(tree)

	return nil
}

func createTreeSha1Hash(tree *objects.TreeObject) {
	nodesContent := strings.Builder{}
	for _, child := range tree.Children {
		nodesContent.WriteString(fmt.Sprintf("%s %s\x00%s", child.Value.Mode, child.Value.Name, child.Value.Hash))
	}
	nodesContentStr := nodesContent.String()

	tree.Content = nodesContentStr
	tree.Size = len(nodesContentStr)

	c := fmt.Sprintf("tree %d\000%s", tree.Size, tree.Content)
	h := sha1.Sum([]byte(c))
	tree.Value.Hash = string(h[:])
}
