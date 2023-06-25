package command

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/mesuutt/git-mirror/internal/config"
	"github.com/mesuutt/git-mirror/pkg/commit"
	"github.com/mesuutt/git-mirror/pkg/git"
	"github.com/mesuutt/git-mirror/pkg/repo"
)

var AddCmd = &cli.Command{
	Name:  "add",
	Usage: "add stats of latest commit",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "path",
			Value: ".",
			Usage: "git repo directory path",
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Value: false,
			Usage: "print changes to stdout instead create commit",
		},
	},
	Action: AddCmdAction,
}

// AddCmdAction add stats of the latest commit to mirror repository
func AddCmdAction(ctx *cli.Context) error {
	// TODO: ignore merge commits
	statRepoPath := ctx.String("repo")
	runPath := ctx.String("path")
	dryRun := ctx.Bool("dry-run")

	if err := git.ValidateRepo(runPath); err != nil {
		return err
	}

	if err := git.ValidateRepo(statRepoPath); err != nil {
		return err
	}

	out, err := git.LastCommitStatsDiff(runPath)
	if err != nil {
		return err
	}

	parser := commit.NewParser()

	if ctx.String("whitelist") != "" {
		parser = parser.WithWhitelist(strings.Split(ctx.String("whitelist"), ","))
	}

	stats, err := parser.Parse(bytes.NewReader(out))
	if err != nil {
		return fmt.Errorf("diff output parse failed with error: `%w`", err)
	}

	if len(stats) == 0 {
		return nil
	}

	conf, err := config.ReadConfig(filepath.Join(statRepoPath, "config.toml"))
	if err != nil {
		// print error and continue with defaults when cannot read config file.
		fmt.Printf("config file read err. Running with defaults. err: %v\n", err)
		conf = config.Default()
	}

	commitGen := commit.NewDiffGenerator(conf)

	var repoImpl repo.Repo
	if dryRun {
		repoImpl = repo.NewFakeRepo()
	} else {
		repoImpl = repo.NewFsRepo(statRepoPath)
	}

	// TODO: ignore already added commit
	// if user run add multiple times without new commit, it should add only one commit to repo
	// maybe we can add commit hast to commit message, and check it at next commit
	diff, err := commitGen.GenDiff(stats, commit.Meta{Time: time.Now()})
	if err != nil {
		return err
	}

	if err := repoImpl.AddStats(diff); err != nil {
		return err
	}

	if err := repoImpl.AddAndCommit("update"); err != nil {
		return err
	}

	fmt.Println("commit stats added to git-mirror repository")
	return nil
}
