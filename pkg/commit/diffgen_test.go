package commit

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mesuutt/git-mirror/pkg/config"
)

func TestGenerateDiffWithDefaultsWhenConfigNotGiven(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: ".go"}

	gen := NewDiffGenerator(config.Default())
	// when
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: time.Now()})

	// then
	assert.Nil(t, err)
	change := diff.Changes[0]
	assert.Equal(t, "log.go", change.Filename)
	assert.Equal(t, filepath.Join(strings.Split(time.Now().Format("2006-01-02"), "-")...), change.Dir)
	assert.Equal(t, fmt.Sprintf("%d insertion(s), %d deletion(s)", stat.Insert, stat.Delete), change.Text)
}

func TestGenerateDiffFilenameShouldNotChangedWhenNotHasDot(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: "Makefile"}

	gen := NewDiffGenerator(config.Default())
	// when
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: time.Now()})

	// then
	assert.Nil(t, err)
	assert.Equal(t, "Makefile", diff.Changes[0].Filename)
}

func TestGenerateDiffWithCustomLogTemplate(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: ".go"}

	gen := NewDiffGenerator(&config.Config{
		Commit: config.CommitConfig{
			Template: "insert: {{.InsertCount}}, delete: {{.DeleteCount}}, fileType: {{.FileType}}, HM: {{.HM}}, Hour: {{.Hour}}, Minute: {{.Minute}}",
		},
	})

	// when
	commitTime := time.Now()
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: commitTime})

	// then
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf("insert: %d, delete: %d, fileType: %s, HM: %s, Hour: %s, Minute: %s", stat.Insert, stat.Delete, strings.ReplaceAll(stat.Ext, ".", ""), commitTime.Format("15:04"), commitTime.Format("15"), commitTime.Format("04")), diff.Changes[0].Text)
}

func TestGenerateDiffWithFileTypeOverwrites(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: ".mod"}
	stat2 := FileStat{Insert: 1, Delete: 2, Ext: ".go"}
	stat3 := FileStat{Insert: 2, Delete: 3, Ext: "Dotless1"}

	gen := NewDiffGenerator(&config.Config{
		Overwrites: map[string]map[string]string{
			"default": {
				"mod":      "go",
				"Dotless1": "Dotless2",
			},
		},
	})

	// when
	diff, err := gen.GenDiff([]FileStat{stat, stat2, stat3}, CommitInfo{Time: time.Now()})

	// then
	assert.Nil(t, err)
	assert.Equal(t, 2, len(diff.Changes))
	assert.Equal(t, "log.go", diff.Changes[0].Filename)
	assert.Equal(t, 2, diff.Changes[0].Insertion)
	assert.Equal(t, 4, diff.Changes[0].Deletion)
	assert.Equal(t, "Dotless2", diff.Changes[1].Filename)
	assert.Equal(t, 2, diff.Changes[1].Insertion)
	assert.Equal(t, 3, diff.Changes[1].Deletion)
}

func TestGenerateDiffWithCodeTemplates(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: ".go"}

	gen := NewDiffGenerator(&config.Config{
		Commit: config.CommitConfig{Template: "my msg"},
		Templates: map[string]string{
			"go": fmt.Sprintf("fmt.Println(\"{{.HM}},{{.Hour}},{{.Minute}},{{.Message}}\")"),
		},
	})

	// when
	commitTime := time.Now()
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: commitTime})

	// then
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("fmt.Println(\"%s,%s,%s,%s\")",
		commitTime.Format("15:04"),
		commitTime.Format("15"),
		commitTime.Format("04"),
		"my msg",
	), diff.Changes[0].Text)
}
