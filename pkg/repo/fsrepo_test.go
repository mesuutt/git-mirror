package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mesuutt/git-mirror/pkg/commit"
)

func TestRepo_AddChange(t *testing.T) {
	repoPath, _ := os.MkdirTemp("", "")
	t.Cleanup(func() { os.RemoveAll(repoPath) })

	repo := NewFsRepo(repoPath)

	tests := []struct {
		name    string
		changes []commit.Change
	}{
		{
			name: "one commit",
			changes: []commit.Change{
				{Insertion: 1, Deletion: 2, Filename: "log.go", Dir: "2023/01/01", Text: "1 insertion(s), 2 deletion(s)"},
			},
		},
		{
			name: "two commit",
			changes: []commit.Change{
				{Insertion: 3, Deletion: 4, Filename: "log.rs", Dir: "2023/01/02", Text: "3 insertion(s), 4 deletion(s)"},
				{Insertion: 2, Deletion: 3, Filename: "log.go", Dir: "2023/01/02", Text: "2 insertion(s), 3 deletion(s)"},
			},
		},
		{
			name:    "file without extension",
			changes: []commit.Change{{Insertion: 1, Deletion: 2, Filename: "Makefile", Dir: "2023/01/03", Text: "1 insertion(s), 2 deletion(s)"}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := repo.AddStats(commit.Diff{Changes: tt.changes})

			if err != nil {
				t.Errorf("stat add failed: %v", err)
			}

			for _, change := range tt.changes {
				path := filepath.Join(repoPath, change.Dir, change.Filename)
				if _, err := os.Stat(path); err != nil {
					t.Fatalf("expected stat file not found in repo: %s", path)
				}

				content, _ := os.ReadFile(path)

				expectedContent := fmt.Sprintf("%s\n", change.Text)
				if string(content) != expectedContent {
					t.Fatalf("file content not mathed. want: `%s`, got: `%s`", expectedContent, string(content))
				}
			}

		})
	}
}
