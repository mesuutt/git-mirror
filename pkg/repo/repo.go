package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mesuutt/git-mirror/pkg/git"
	"github.com/mesuutt/git-mirror/pkg/parser"
)

type Repo struct {
	path string
}

func NewRepo(path string) Repo {
	return Repo{path: path}
}

// AddStats writes diff stats to related files in mirror repo
func (r Repo) AddStats(stats ...parser.FileStat) error {
	dayParts := strings.Split(time.Now().Format("2006-01-02"), "-")

	dir := filepath.Join(r.path, dayParts[0], dayParts[1], dayParts[2])

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
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

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
