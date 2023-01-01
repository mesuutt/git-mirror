package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCmd_ValidateRunRepo(t *testing.T) {
	runPath, _ := os.MkdirTemp("", "")
	t.Cleanup(func() { os.RemoveAll(runPath) })

	app := App()

	err := app.Run([]string{"executable-name", "add", "--path=" + runPath})
	assert.ErrorContains(t, err, "is not a git repository")
}

func TestAddCmd_ValidateStatsRepo(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "")
	t.Cleanup(func() { os.RemoveAll(tempDir) })
	runRepoPath := filepath.Join(tempDir, "run-repo")
	createRepo(t, runRepoPath)
	notExistingPath := filepath.Join(tempDir, "stats-repo")

	app := App()

	err := app.Run([]string{"executable-name", "--repo=" + notExistingPath, "add", "--path=" + runRepoPath})
	assert.ErrorContains(t, err, "is not a git repository")
}

func createRepo(t *testing.T, repoPath string) {
	t.Helper()

	if err := os.MkdirAll(repoPath, os.ModePerm); err != nil {
		t.Fatalf("test git repo create failed: %v", err)
	}

	_, err := exec.Command("git", "-C", repoPath, "init").Output()
	if err != nil {
		t.Fatalf("git repo init failed: %s, `%v`", repoPath, err)
	}
}

// get the latest commit
