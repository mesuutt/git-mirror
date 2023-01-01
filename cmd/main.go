package main

import (
	"log"
	"os"

	"gitmirror/commands"
)

func main() {
	if err := commands.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
