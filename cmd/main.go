package main

import (
	"log"
	"os"

	"gitmirror/command"
)

func main() {
	if err := command.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
