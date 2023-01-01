package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("user home directory path getting failed: %v", err)
	}

	return &cli.App{
		Usage: "git activity mirror",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:  "repo",
				Value: fmt.Sprintf(filepath.Join(homeDir, ".git-mirror")),
				Usage: "git repo directory path of the mirror repo",
			},
			&cli.StringFlag{
				Name:    "whitelist",
				Value:   "",
				Usage:   "comma seperated file extensions to create stats. eg: py,go,sh,Makefile",
				EnvVars: []string{"GIT_MIRROR_FILE_TYPE_WHITELIST"},
			},
		},
		Commands: []*cli.Command{
			{
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
			},
			AddCmd,
		},
	}
}
