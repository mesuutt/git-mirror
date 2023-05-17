package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mesuutt/git-mirror/pkg/commit"
	"github.com/mesuutt/git-mirror/pkg/git"
)

type Repo struct {
	path      string
	commitGen commit.Generator
}

func NewRepo(path string, commitGen commit.Generator) Repo {
	return Repo{path: path, commitGen: commitGen}
}

// AddStats writes diff stats to related files in mirror repo
func (r Repo) AddStats(stats ...commit.FileStat) error {
	dayParts := strings.Split(time.Now().Format("2006-01-02"), "-")

	dir := filepath.Join(r.path, dayParts[0], dayParts[1], dayParts[2])

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("directory could not created: `%s`. error: %w", dir, err)
	}

	for i := range stats {
		stat := &stats[i]

		commit := r.commitGen.GenCommit(stat)
		filePath := filepath.Join(dir, commit.Filename)

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("file could not be opened for appending changes: `%s`. error: %w", filePath, err)
		}

		_, err = f.WriteString(commit.Message + "\n")
		if err != nil {
			return fmt.Errorf("changes could not write to file: `%s`. error: %w", filePath, err)
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

// AddAndCommit add latest generated stats to git and commit
func (r Repo) AddAndCommit(msg string) error {
	if err := git.AddChanges(r.path); err != nil {
		return err
	}

	return git.Commit(r.path, msg)
}
