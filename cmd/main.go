package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

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
				Usage:   "comma seperated file extensions to create stats. eg: go,rs,sh,Makefile",
				EnvVars: []string{"GIT_MIRROR_FILE_TYPE_WHITELIST"},
			},
		},
		Commands: []*cli.Command{
			InstallCmd,
			AddCmd,
		},
	}
}
