package command

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"gitmirror"
	"gitmirror/git"
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
	},
	Action: AddCmdAction,
}

// AddCmdAction add stats of the latest commit to mirror repository
func AddCmdAction(ctx *cli.Context) error {
	// TODO: ignore merge commits
	statRepoPath := ctx.String("repo")
	runPath := ctx.String("path")

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

	whitelist := strings.Split(ctx.String("whitelist"), ",")
	parser := gitmirror.NewParser().WithWhitelist(whitelist)

	stats, err := parser.Parse(bytes.NewReader(out))
	if err != nil {
		return fmt.Errorf("diff output parse failed with error: `%v`", err)
	}

	// TODO: ignore already added commit
	// if user run add multiple times without new commit, it should add only one commit to repo
	repo := gitmirror.NewRepo(statRepoPath)

	if err := repo.AddStats(afero.NewMemMapFs(), stats...); err != nil {
		return err
	}

	if err := repo.Commit("update"); err != nil {
		return err
	}

	fmt.Println("commit stats added to git-mirror repository")
	return nil
}
