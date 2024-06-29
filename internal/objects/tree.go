package objects

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type ObjectModes string

const (
	RegularFileMode    = "100644"
	ExecutableFileMode = "100755"
	SymbolicLinkMode   = "120000"
	DirectoryMode      = "040000"
)

type TreeObjectNode struct {
	Type ObjectType `json:"type"`
	Mode string     `json:"mode"`
	Name string     `json:"name"`
	Hash string     `json:"hash"`
}

type TreeObject struct {
	CommonObject
	Children []TreeObjectNode `json:"children"`
}

func ParseTreeObject(object CommonObject) Object {
	nodes := make([]TreeObjectNode, 0)
	last := TreeObjectNode{}

	idx := 0
	s := strings.Builder{}

	for idx < len(object.Content) {
		// <mode> <name>\0<20_byte_sha>
		// <mode> <name>\0<20_byte_sha>
		str := object.Content[idx]
		if string(str) == " " {
			last.Mode = s.String()
			s.Reset()

			// we cover only this mode for now
			if last.Mode == DirectoryMode {
				last.Type = Tree
			} else {
				last.Type = Blob
			}

			nullIdx := strings.Index(object.Content[idx:], "\x00")

			name := object.Content[idx+1 : idx+nullIdx]
			last.Name = name

			idx += nullIdx + 1

			reader := strings.NewReader(object.Content[idx:])

			hash := make([]byte, 20)
			_, _ = io.ReadAtLeast(reader, hash, 20)

			last.Hash = hex.EncodeToString(hash[:])

			// 20 bytes (hash) + 1 byte (jump to next node)
			idx += 20

			nodes = append(nodes, last)
			last = TreeObjectNode{}
			continue
		}

		_, err := s.WriteString(string(str))
		if err != nil {
			panic(err)
		}
		idx++
	}

	return TreeObject{
		CommonObject: object,
		Children:     nodes,
	}
}

func ParseTreeObjectContent(nodeStr string) TreeObjectNode {
	node := TreeObjectNode{}

	idx := strings.Index(nodeStr, " ")
	node.Mode = nodeStr[:idx]

	nodeStr = nodeStr[idx+1:]
	idx = strings.Index(nodeStr, " ")
	node.Type = ObjectType(nodeStr[:idx])

	nodeStr = nodeStr[idx+1:]
	idx = strings.Index(nodeStr, "\x00")
	node.Hash = nodeStr[:idx]

	return node
}

func (t TreeObject) String() string {
	b := strings.Builder{}
	for _, content := range t.Children {
		b.WriteString(fmt.Sprintf("%s %s %s    %s\n", content.Mode, content.Type, content.Hash, content.Name))
	}
	return b.String()
}
