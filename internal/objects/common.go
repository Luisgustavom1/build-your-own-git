package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
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
	Type ObjectType `json:"type"`
	Size int        `json:"size"`
	Hash string     `json:"hash"`
	// TODO: review this, maybe Data and Content can be merged
	Data    string `json:"data"`
	Content string `json:"content"`
}

func NewCommonObjectFromHash(hash string) CommonObject {
	content, err := RepoGetObject(hash)
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
	c := CommonObject{Content: blob}

	typeIdx := strings.IndexByte(blob, ' ')
	c.Type = ObjectType(blob[:typeIdx])

	idx := strings.IndexByte(blob, '\x00')

	size, err := strconv.Atoi(blob[typeIdx+1 : idx])
	if err != nil {
		panic(err)
	}

	c.Size = size
	c.Data = blob[idx+1:]

	return c
}

func RepoGetObject(name string) ([]byte, error) {
	return os.ReadFile(path.Join(".git/objects", name[:2], name[2:]))
}

func RepoCheckObjectId(name string) bool {
	_, err := os.Stat(path.Join(".git/objects", name[:2], name[2:]))
	return err == nil
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

func SaveObject(obj CommonObject) error {
	sha1_hash := obj.Hash
	objectPath := path.Join(".git/objects", string(sha1_hash[:2]))
	objectFile := sha1_hash[2:]

	err := os.MkdirAll(objectPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory -> %s\n", err)
	}

	compressedBlob := bytes.Buffer{}
	w := zlib.NewWriter(&compressedBlob)
	w.Write([]byte(obj.Content))
	w.Close()

	err = os.WriteFile(path.Join(objectPath, objectFile), compressedBlob.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Error writing file -> %s\n", err)
	}

	return nil
}

func CreateObjectHash(blob []byte) string {
	h := sha1.Sum(blob)
	hash := hex.EncodeToString(h[:])
	return hash
}
