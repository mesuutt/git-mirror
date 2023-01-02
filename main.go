package main

import (
	"log"
	"os"

	"github.com/mesuutt/git-mirror/internal/command"
)

func main() {
	if err := command.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
