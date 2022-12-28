package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func InstallHookCmd(cCtx *cli.Context) error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		log.Fatalf("you need to run this command in a git repository.")
	}

	hookDir := filepath.Join(".git", "hooks")
	if _, err := os.Stat(hookDir); os.IsNotExist(err) {
		if err := os.MkdirAll(hookDir, os.ModePerm); err != nil {
			log.Fatalf("%s directory creation failed", hookDir)
		}
	}

	hookFilePath := filepath.Join(hookDir, "post-commit")
	f, err := os.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("%s file could not open for adding hook", hookFilePath)
	}

	defer f.Close()

	// TODO: get app name from args (git-mirror)
	if _, err := f.WriteString("git-mirror add"); err != nil {
		log.Fatalf("hook add failed to post-commit file. error: %v", err)
	}

	fmt.Printf("post-commit hook installed at %s\n", hookFilePath)
	return nil
}
