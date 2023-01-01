package gitmirror

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/afero"

	"gitmirror/git"
)

type Repo struct {
	path string
}

func NewRepo(path string) Repo {
	return Repo{path: path}
}

// AddStats writes diff stats to related files in mirror repo
func (r Repo) AddStats(fs afero.Fs, stats ...FileStat) error {
	now := time.Now()
	dir := filepath.Join(r.path, strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()))

	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("directory could not created: `%s`. error: %w", dir, err)
	}

	for i := range stats {
		stat := &stats[i]
		filename := fmt.Sprintf("log%s", stat.Ext)

		// handle files without extension. eg: Makefile
		if !strings.HasPrefix(stat.Ext, ".") {
			filename = stat.Ext
		}

		filePath := filepath.Join(dir, filename)

		f, err := fs.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("file could not be opened for appending changes: `%s`. error: %w", filePath, err)
		}

		_, err = f.WriteString(fmt.Sprintf("%d insertion(s), %d deletion(s)\n", stat.Insert, stat.Delete))
		if err != nil {
			return fmt.Errorf("changes could not write to file: `%s`. error: %w", filePath, err)
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Commit commits changes
func (r Repo) Commit(msg string) error {
	return git.AddAndCommit(r.path, msg)
}
