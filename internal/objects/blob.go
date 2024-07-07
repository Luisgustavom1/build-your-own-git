package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

type BlobObject struct {
	CommonObject
	Data string
	Hash string
}

func NewBlobObject(data []byte) BlobObject {
	common := CommonObject{Type: Blob, Content: string(data), Size: len(data)}

	blob := fmt.Sprintf("blob %d\000%s", len(data), data)
	hash := CreateBlobHash(blob)

	object := BlobObject{CommonObject: common, Data: blob, Hash: hash}

	return object
}

func (b BlobObject) Save() error {
	sha1_hash := b.Hash
	objectPath := path.Join(".git/objects", string(sha1_hash[:2]))
	objectFile := sha1_hash[2:]

	err := os.MkdirAll(objectPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory -> %s\n", err)
	}

	compressedBlob := bytes.Buffer{}
	w := zlib.NewWriter(&compressedBlob)
	w.Write([]byte(b.Data))
	w.Close()

	err = os.WriteFile(path.Join(objectPath, objectFile), compressedBlob.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Error writing file -> %s\n", err)
	}

	return nil
}

func (b BlobObject) String() string {
	return b.CommonObject.Content
}

func ParseBlobObject(object CommonObject) BlobObject {
	blob := BlobObject{CommonObject: object}
	return blob
}

func CreateBlobHash(blob string) string {
	h := sha1.Sum([]byte(blob))
	hash := hex.EncodeToString(h[:])
	return hash
}
