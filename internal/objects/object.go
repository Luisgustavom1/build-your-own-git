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

type Object struct {
	Type ObjectType `json:"type"`
	Size int        `json:"size"`
	Hash string     `json:"hash"`
	Data string     `json:"data"`
}

func NewObject(ttype ObjectType, size int, data string) Object {
	obj := Object{
		Type: ttype,
		Size: size,
		Data: data,
	}
	obj.GenerateHash()
	return obj
}

func NewObjectFromHash(hash string) Object {
	content, err := RepoGetObject(hash)
	if err != nil {
		panic(fmt.Errorf("Error reading object -> %s\n", err))
	}

	data, err := uncompressObjectContent(content)
	if err != nil {
		panic(fmt.Errorf("Error decompressing object -> %s\n", err))
	}

	return NewObjectFromContent(data)
}

func NewObjectFromContent(content string) Object {
	c := Object{}

	typeIdx := strings.IndexByte(content, ' ')
	c.Type = ObjectType(content[:typeIdx])

	idx := strings.IndexByte(content, '\x00')

	size, err := strconv.Atoi(content[typeIdx+1 : idx])
	if err != nil {
		panic(err)
	}

	c.Size = size
	c.Data = content[idx+1:]

	return c
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

func (obj *Object) Write() error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git\n")
	}

	sha1_hash := obj.Hash
	objectPath := path.Join(".git/objects", string(sha1_hash[:2]))
	objectFile := sha1_hash[2:]

	err := os.MkdirAll(objectPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory -> %s\n", err)
	}

	compressedBlob := bytes.Buffer{}
	w := zlib.NewWriter(&compressedBlob)
	w.Write([]byte(obj.GetContent()))
	w.Close()

	err = os.WriteFile(path.Join(objectPath, objectFile), compressedBlob.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Error writing file -> %s\n", err)
	}

	return nil
}

func (o *Object) GenerateHash() string {
	o.Hash = CreateObjectHash(o.GetContent())
	return o.Hash
}

func (o *Object) GetContent() string {
	return fmt.Sprintf("%s %d\x00%s", o.Type, o.Size, o.Data)
}

func CreateObjectHash(blob string) string {
	h := sha1.Sum([]byte(blob))
	hash := hex.EncodeToString(h[:])
	return hash
}
