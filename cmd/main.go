package main

import (
	"log"
	"os"

	"gitmirror/command"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("user home directory path getting failed: %v", err)
	}

	if err := command.InitCLI(homeDir).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
