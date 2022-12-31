package command

import (
	"context"
	"fmt"
	"path"

	"github.com/spf13/afero"
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
		Before: func(ctx *cli.Context) error {
			ctx.Context = context.WithValue(ctx.Context, "fs", afero.NewOsFs())
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:            "install",
				Usage:           "install post-commit hook for adding stats automatically",
				SkipFlagParsing: true,
				Action:          InstallHookCmd,
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
