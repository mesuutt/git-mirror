package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallHookCmd(t *testing.T) {
	repoPath, _ := os.MkdirTemp("", "")
	t.Cleanup(func() { os.RemoveAll(repoPath) })

	hookDir := filepath.Join(repoPath, ".git", "hooks")
	hookFile := filepath.Join(hookDir, "post-commit")
	_ = os.MkdirAll(hookDir, os.ModePerm)

	// when
	err := App().Run([]string{"executable-name", "install", "--path=" + repoPath})

	// then
	assert.Nil(t, err)

	bytes, err := os.ReadFile(hookFile)
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), "git-mirror add")
}
