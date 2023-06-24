package repo

import (
	"fmt"
	"path"

	"github.com/mesuutt/git-mirror/pkg/commit"
)

type stdoutRepo struct {
}

var _ Repo = (*stdoutRepo)(nil)

func NewFakeRepo() stdoutRepo {
	return stdoutRepo{}
}

func (r stdoutRepo) AddStats(diff *commit.Diff) error {
	for _, change := range diff.Changes {
		filePath := path.Join(change.Dir, change.Filename)
		fmt.Printf("%s\n\t%s\n", filePath, change.Text)
	}

	return nil
}

func (r stdoutRepo) AddAndCommit(msg string) error {
	fmt.Printf("commit message: \n\t%s\n", msg)
	return nil
}
