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
	DirectoryMode      = "40000"
)

type TreeObjectNode struct {
	CommonObject
	Mode string `json:"mode"`
	Name string `json:"name"`
}

type TreeObject struct {
	Value    TreeObjectNode `json:"value"`
	Children []*TreeObject  `json:"children"`
}

type TreeStringOpts struct {
	NameOnly bool
}

func ParseTreeObject(object CommonObject) TreeObject {
	nodes := make([]*TreeObject, 0)
	last := &TreeObject{}

	idx := 0
	s := strings.Builder{}

	for idx < len(object.Data) {
		// <mode> <name>\0<20_byte_sha>
		// <mode> <name>\0<20_byte_sha>
		str := object.Data[idx]
		if string(str) == " " {
			last.Value.Mode = s.String()
			s.Reset()

			// we cover only this mode for now
			if last.Value.Mode == DirectoryMode {
				last.Value.Type = Tree
			} else {
				last.Value.Type = Blob
			}

			nullIdx := strings.IndexByte(object.Data[idx:], '\x00')

			name := object.Data[idx+1 : idx+nullIdx]
			last.Value.Name = name

			idx += nullIdx + 1

			reader := strings.NewReader(object.Data[idx:])

			hash := make([]byte, 20)
			_, _ = io.ReadAtLeast(reader, hash, 20)

			last.Value.Hash = hex.EncodeToString(hash[:])

			// 20 bytes (hash)
			idx += 20

			nodes = append(nodes, last)
			last = &TreeObject{}
			continue
		}

		err := s.WriteByte(str)
		if err != nil {
			panic(err)
		}
		idx++
	}

	return TreeObject{
		Value: TreeObjectNode{
			CommonObject: object,
		},
		Children: nodes,
	}
}

func (t TreeObject) String(opts TreeStringOpts) string {
	b := strings.Builder{}
	for _, content := range t.Children {
		if opts.NameOnly {
			b.WriteString(fmt.Sprintf("%s\n", content.Value.Name))
		} else {
			b.WriteString(fmt.Sprintf("%06s %s %s    %s\n", content.Value.Mode, content.Value.Type, content.Value.Hash, content.Value.Name))
		}
	}
	return b.String()
}
