package main

import (
	"log"
	"os"

	"github.com/mesuutt/git-mirror/commands"
)

func main() {
	if err := commands.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
