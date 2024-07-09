package objects

type BlobObject struct {
	Object
}

func NewBlobObject(object Object) BlobObject {
	blob := BlobObject{Object: object}
	return blob
}

func NewBlobObjectFromData(data []byte) BlobObject {
	object := NewObject(Blob, len(data), string(data))

	return BlobObject{Object: object}
}

func (b BlobObject) String() string {
	return b.Object.Data
}
