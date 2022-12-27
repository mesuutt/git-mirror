package main

import (
	"fmt"
	"path"

	"github.com/urfave/cli/v2"
)

func InitCLI(homeDir string) *cli.App {
	return &cli.App{
		Usage: "git activity mirror",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "repo",
				Value: fmt.Sprintf(path.Join(homeDir, ".git-mirror")),
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
				Action: func(cCtx *cli.Context) error {
					// TODO
					fmt.Println("post-commit installed at .git/hooks/post-commit")
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "add stats of latest commit",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "path",
						Value: ".",
						Usage: "git repo directory path",
					},
				},
				Action: AddCmd,
			},
		},
	}
}
