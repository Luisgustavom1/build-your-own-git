package objects

import (
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
	Type    ObjectType `json:"type"`
	Size    int        `json:"size"`
	Content string     `json:"content"`
}

type Object interface {
	String() string
}

func ParseCommonObject(blob string) CommonObject {
	c := CommonObject{}

	typeIdx := strings.Index(blob, " ")
	c.Type = ObjectType(blob[:typeIdx])

	blob = blob[typeIdx+1:]
	idx := strings.Index(blob, "\x00")

	size, err := strconv.Atoi(blob[:idx])
	if err != nil {
		panic(err)
	}
	c.Size = size

	c.Content = blob[idx+1:]

	return c
}
