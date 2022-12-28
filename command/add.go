package command

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"gitmirror"
)

func AddCmd(cCtx *cli.Context) error {
	// TODO: ignore merge commits
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
		log.Fatalf("`git -C %s diff --numstat HEAD~1` command execution failed. error %v", repoPath, err)
	}

	// TODO: ignore already added commit
	// if user run add multiple times without new commit, it should add only one commit to repo
	repo := gitmirror.NewRepo(statRepoPath)
	parser := gitmirror.NewParser(nil)
	stats, err := parser.Parse(bytes.NewReader(out))
	if err != nil {
		log.Fatal(err)
	}

	if err := repo.AddStats(fs, stats...); err != nil {
		return err
	}

	fmt.Println("commit stats added to git-mirror repository")
	return nil
}
