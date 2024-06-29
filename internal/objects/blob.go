package objects

type BlobObject struct {
	CommonObject
	data string
}

func ParseBlobObject(object CommonObject) Object {
	blob := BlobObject{CommonObject: object}
	blob.data = object.Content
	return blob
}

func (b BlobObject) String() string {
	return b.data
}
