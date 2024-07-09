package objects

import (
	"fmt"
)

type BlobObject struct {
	CommonObject
}

func ParseBlobObject(object CommonObject) BlobObject {
	blob := BlobObject{CommonObject: object}
	return blob
}

func NewBlobObject(data []byte) BlobObject {
	blob := fmt.Sprintf("blob %d\000%s", len(data), data)
	hash := CreateObjectHash([]byte(blob))

	common := CommonObject{Type: Blob, Content: blob, Size: len(data), Data: string(data), Hash: hash}
	object := BlobObject{CommonObject: common}

	return object
}

func (b BlobObject) String() string {
	return b.CommonObject.Data
}
