package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func AddCmd(cCtx *cli.Context) error {
	fs := afero.NewOsFs()
	statRepoPath := cCtx.String("repo")
	repoPath := cCtx.String("path")

	if _, err := fs.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
		log.Fatalf("%s is not a git repository.", repoPath)
	}

	if fileInfo, err := os.Stat(statRepoPath); os.IsNotExist(err) {
		log.Fatalf("git mirror repo not exists. Please create a repo on Github/Gitlab etc and clone it to the %s path", statRepoPath)
	} else if !fileInfo.IsDir() {
		log.Fatalf("%s is not a directory. given path should be a git repo", statRepoPath)
	}

	out, err := exec.Command("git", "-C", repoPath, "diff", "--numstat", "HEAD~1").Output()
	if err != nil {
		log.Fatal(err)
	}

	repo := NewRepo(statRepoPath)
	parser := NewParser(nil)
	stats, err := parser.Parse(bytes.NewReader(out))
	if err != nil {
		log.Fatal(err)
	}

	return repo.AddStats(fs, stats...)
}
