package objects

import (
	"os"
	"path"
)

func RepoGetObject(name string) ([]byte, error) {
	return os.ReadFile(path.Join(".git/objects", name[:2], name[2:]))
}

func RepoCheckObjectId(name string) bool {
	_, err := os.Stat(path.Join(".git/objects", name[:2], name[2:]))
	return err == nil
}
