package objects

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type BlobObject struct {
	CommonObject
}

func NewBlobObject(data []byte) BlobObject {
	blob := fmt.Sprintf("blob %d\000%s", len(data), data)
	hash := createBlobHash(blob)

	common := CommonObject{Type: Blob, Content: blob, Size: len(data), Data: string(data), Hash: hash}
	object := BlobObject{CommonObject: common}

	return object
}

func (b BlobObject) String() string {
	return b.CommonObject.Data
}

func ParseBlobObject(object CommonObject) BlobObject {
	blob := BlobObject{CommonObject: object}
	return blob
}

func createBlobHash(blob string) string {
	h := sha1.Sum([]byte(blob))
	hash := hex.EncodeToString(h[:])
	return hash
}
