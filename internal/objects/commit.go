package objects

import (
	"fmt"
	"time"
)

type CommitObject struct {
	Object
}

const GIT_AUTHOR_EMAIL = "author@example.com"
const GIT_AUTHOR_NAME = "Author Name"

var Now = time.Now()

func NewCommitObject(object Object) CommitObject {
	commit := CommitObject{Object: object}
	return commit
}

func NewCommitFromData(data string) CommitObject {
	obj := NewObject(Commit, len(data), data)

	return CommitObject{
		Object: obj,
	}
}

func (c CommitObject) String() string {
	return c.Object.Data
}

func GenerateCommitData(treeHash, parentHash, message string) string {
	seconds := Now.Unix()
	z, _ := Now.Zone()

	headers := fmt.Sprintf("tree %s", treeHash)
	if parentHash != "" {
		headers += fmt.Sprintf("\nparent %s", parentHash)
	}

	author := fmt.Sprintf("author %s <%s> %d %s", GIT_AUTHOR_NAME, GIT_AUTHOR_EMAIL, seconds, z)
	committer := fmt.Sprintf("committer %s <%s> %d %s", GIT_AUTHOR_NAME, GIT_AUTHOR_EMAIL, seconds, z)
	data := fmt.Sprintf("%s\n%s\n%s\n\n%s\n", headers, author, committer, message)
	return data
}
