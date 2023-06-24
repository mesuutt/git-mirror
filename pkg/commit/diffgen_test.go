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

	gen := NewDiffGenerator(nil)
	// when
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: time.Now()})

	// then
	assert.Nil(t, err)
	change := diff.Changes[0]
	assert.Equal(t, "log.go", change.Filename)
	assert.Equal(t, filepath.Join(strings.Split(time.Now().Format("2006-01-02"), "-")...), change.Dir)
	assert.Equal(t, fmt.Sprintf(defaultLogMsgFormat, stat.Insert, stat.Delete), change.Text)
}

func TestGenerateDiffFilenameShouldNotChangedWhenNotHasDot(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: "Makefile"}

	gen := NewDiffGenerator(nil)
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
			Template: "insert: {{.InsertCount}}, delete: {{.DeleteCount}}, ext: {{.Ext}}, HM: {{.HM}}, Hour: {{.Hour}}, Minute: {{.Minute}}",
		},
	})

	// when
	commitTime := time.Now()
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: commitTime})

	// then
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf("insert: %d, delete: %d, ext: %s, HM: %s, Hour: %s, Minute: %s", stat.Insert, stat.Delete, strings.ReplaceAll(stat.Ext, ".", ""), commitTime.Format("15:04"), commitTime.Format("15"), commitTime.Format("04")), diff.Changes[0].Text)
}

func TestGenerateDiffWithFileExtensionFromConfig(t *testing.T) {
	// given
	stat := FileStat{Insert: 1, Delete: 2, Ext: ".mod"}

	gen := NewDiffGenerator(&config.Config{
		Overwrites: map[string]map[string]string{
			"default": {
				"mod": "go",
			},
		},
	})

	// when
	diff, err := gen.GenDiff([]FileStat{stat}, CommitInfo{Time: time.Now()})

	// then
	assert.Nil(t, err)
	assert.Equal(t, "log.go", diff.Changes[0].Filename)
}
