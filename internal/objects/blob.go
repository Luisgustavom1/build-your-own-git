package objects

type BlobObject struct {
	CommonObject
	data string
}

func ParseBlobObject(object CommonObject) BlobObject {
	blob := BlobObject{CommonObject: object}
	blob.data = object.Content
	return blob
}

func BlobObjToString(b BlobObject) string {
	return b.data
}
