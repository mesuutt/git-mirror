package repo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mesuutt/git-mirror/pkg/commit"
	"github.com/mesuutt/git-mirror/pkg/git"
)

type fsRepo struct {
	path string
}

var _ Repo = (*fsRepo)(nil)

func NewFsRepo(path string) fsRepo {
	return fsRepo{path: path}
}

// AddStats writes diff stats to related files in mirror repo
func (r fsRepo) AddStats(diff commit.Diff) error {
	for i := range diff.Changes {
		change := &diff.Changes[i]
		absDirPath := filepath.Join(r.path, change.Dir)

		if err := os.MkdirAll(absDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("directory could not created: `%s`. error: %w", absDirPath, err)
		}

		filePath := filepath.Join(absDirPath, change.Filename)

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("file could not be opened for appending changes: `%s`. error: %w", filePath, err)
		}

		_, err = f.WriteString(change.Text + "\n")
		if err != nil {
			return fmt.Errorf("changes could not write to file: `%s`. error: %w", filePath, err)
		}

		return f.Close()
	}

	return nil
}

// AddAndCommit add latest generated stats to git and commit
func (r fsRepo) AddAndCommit(msg string) error {
	if err := git.AddChanges(r.path); err != nil {
		return err
	}

	return git.Commit(r.path, msg)

}
