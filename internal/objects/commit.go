package objects

import (
	"fmt"
	"time"
)

type CommitObject struct {
	CommonObject
	TreeHash   string `json:"tree_hash"`
	ParentHash string `json:"parent_hash"`
	Message    string `json:"message"`
}

const GIT_AUTHOR_EMAIL = "author@example.com"
const GIT_AUTHOR_NAME = "Author Name"

var Now = time.Now()

func ParseCommitObject(object CommonObject) CommitObject {
	commit := CommitObject{CommonObject: object}
	return commit
}

func NewCommitObject(treeHash, parentHash, message string) CommitObject {
	seconds := Now.Unix()
	z, _ := Now.Zone()

	headers := fmt.Sprintf("tree %s", treeHash)
	var parent string
	if parentHash != "" {
		parent = fmt.Sprintf("parent %s\n", parentHash)
	}
	author := fmt.Sprintf("author %s <%s> %d %s", GIT_AUTHOR_NAME, GIT_AUTHOR_EMAIL, seconds, z)
	committer := fmt.Sprintf("committer %s <%s> %d %s", GIT_AUTHOR_NAME, GIT_AUTHOR_EMAIL, seconds, z)
	data := fmt.Sprintf("%s\n%s%s\n%s\n\n%s\n", headers, parent, author, committer, message)
	content := fmt.Sprintf("commit %d\x00%s", len(data), data)

	return CommitObject{
		CommonObject: CommonObject{
			Type:    Commit,
			Data:    data,
			Content: content,
			Size:    len(content),
			Hash:    CreateObjectHash([]byte(content)),
		},
		TreeHash:   treeHash,
		ParentHash: parentHash,
		Message:    message,
	}
}

func (c CommitObject) String() string {
	return c.CommonObject.Data
}
