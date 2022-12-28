package gitmirror

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestRepo_AddChange(t *testing.T) {
	appFS := afero.NewMemMapFs()
	repo := NewRepo("/repo")

	now := time.Now()
	tests := []struct {
		name    string
		stats   []FileStat
		path    string
		content string
	}{
		{
			name:    "one commit",
			stats:   []FileStat{{Insert: 1, Delete: 2, Ext: ".go"}},
			path:    fmt.Sprintf("/repo/%d/%d/%d/log.go", now.Year(), now.Month(), now.Day()),
			content: "1 insertion(s), 2 deletion(s)\n",
		},
		{
			name:    "same file type, two commit",
			stats:   []FileStat{{Insert: 1, Delete: 2, Ext: ".rs"}, {Insert: 2, Delete: 2, Ext: ".rs"}},
			path:    fmt.Sprintf("/repo/%d/%d/%d/log.rs", now.Year(), now.Month(), now.Day()),
			content: "1 insertion(s), 2 deletion(s)\n2 insertion(s), 2 deletion(s)\n",
		},
		{
			name:    "file without extension",
			stats:   []FileStat{{Insert: 1, Delete: 2, Ext: "Makefile"}},
			path:    fmt.Sprintf("/repo/%d/%d/%d/Makefile", now.Year(), now.Month(), now.Day()),
			content: "1 insertion(s), 2 deletion(s)\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := repo.AddStats(appFS, tt.stats...)
			if err != nil {
				t.Errorf("stat add failed: %v", err)
			}
			exists, _ := afero.Exists(appFS, tt.path)

			if !exists {
				t.Fatalf("expected stat file not found in desired path: %s", tt.path)
			}

			content, _ := afero.ReadFile(appFS, tt.path)

			if string(content) != tt.content {
				t.Fatalf("file content not mathed. want: `%s`, got: `%s`", tt.content, string(content))
			}
		})
	}
}
