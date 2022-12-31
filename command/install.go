package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

// InstallHookCmd add git hook to .git/hooks/post-commit file
func InstallHookCmd(c *cli.Context) error {
	fs := c.Context.Value("fs").(afero.Fs)
	if _, err := fs.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("you need to run this command in a git repository: %v", err)
	}

	hookDir := filepath.Join(".git", "hooks")
	if _, err := fs.Stat(hookDir); os.IsNotExist(err) {
		if err := fs.MkdirAll(hookDir, os.ModeDir); err != nil {
			return fmt.Errorf(hookDir+" directory creation failed: %v", err)
		}
	}

	// TODO: in test check file is executable
	hookFilePath := filepath.Join(hookDir, "post-commit")
	f, err := fs.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
