package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
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
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("you need to run this command in a git repository: %v", err)
	}

	hookDir := filepath.Join(path, ".git", "hooks")
	if _, err := os.Stat(hookDir); os.IsNotExist(err) {
		if err := os.MkdirAll(hookDir, os.ModeDir); err != nil {
			return fmt.Errorf(hookDir+" directory creation failed: %v", err)
		}
	}

	// TODO: in test check file is executable
	hookFilePath := filepath.Join(hookDir, "post-commit")
	f, err := os.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf(hookFilePath+" file could not open for adding hook: %v", err)
	}

	defer f.Close()

	// TODO: get app name from args (git-mirror)
	if _, err := f.WriteString("git-mirror add"); err != nil {
		fmt.Errorf("hook add failed to post-commit file. error: %v", err)
	}

	fmt.Printf("post-commit hook installed at %s\n", hookFilePath)
	return nil
}
