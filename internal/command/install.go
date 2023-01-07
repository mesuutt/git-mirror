package command

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/mesuutt/git-mirror/pkg/git"
)

var InstallCmd = &cli.Command{
	Name:  "install",
	Usage: "install post-commit hook for adding stats automatically",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:  "path",
			Value: ".",
			Usage: "git repo to install post-commit hook",
		},
	},
	Action: InstallHookCmd,
}

// InstallHookCmd add git hook to .git/hooks/post-commit file
func InstallHookCmd(c *cli.Context) error {
	path := c.Path("path")
	if err := git.ValidateRepo(path); err != nil {
		return err
	}

	if err := git.CreateHookDir(path); err != nil {
		return err
	}

	if err := git.AddHook(path, "post-commit", "git-mirror add"); err != nil {
		return err
	}

	fmt.Println("post-commit hook added")
	return nil
}
