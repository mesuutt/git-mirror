package repo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mesuutt/git-mirror/pkg/commit"
)

func TestRepo_AddChange(t *testing.T) {
	repoPath, _ := os.MkdirTemp("", "")
	t.Cleanup(func() { os.RemoveAll(repoPath) })

	commitGen := commit.NewDiffGenerator(filepath.Join(repoPath, "config.toml"))

	repo := NewFsRepo(repoPath, commitGen)

	now := time.Now()
	tests := []struct {
		name     string
		stats    []commit.FileStat
		fileName string
		content  string
	}{
		{
			name:     "one commit",
			stats:    []commit.FileStat{{Insert: 1, Delete: 2, Ext: ".go"}},
			fileName: "log.go",
			content:  "1 insertion(s), 2 deletion(s)\n",
		},
		{
			name:     "same file type, two commit",
			stats:    []commit.FileStat{{Insert: 1, Delete: 2, Ext: ".rs"}, {Insert: 2, Delete: 2, Ext: ".rs"}},
			fileName: "log.rs",
			content:  "1 insertion(s), 2 deletion(s)\n2 insertion(s), 2 deletion(s)\n",
		},
		{
			name:     "file without extension",
			stats:    []commit.FileStat{{Insert: 1, Delete: 2, Ext: "Makefile"}},
			fileName: "Makefile",
			content:  "1 insertion(s), 2 deletion(s)\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			dayParts := strings.Split(now.Format("2006-01-02"), "-")
			path := filepath.Join(repoPath, dayParts[0], dayParts[1], dayParts[2], tt.fileName)

			err := repo.AddStats(tt.stats...)

			if err != nil {
				t.Errorf("stat add failed: %v", err)
			}

			if _, err := os.Stat(path); err != nil {
				t.Fatalf("expected stat file not found in repo: %s", path)
			}

			content, _ := os.ReadFile(path)

			if string(content) != tt.content {
				t.Fatalf("file content not mathed. want: `%s`, got: `%s`", tt.content, string(content))
			}
		})
	}
}
