package main

import (
	"log"
	"os"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("user home directory path getting failed: %v", err)
	}

	if err := InitCLI(homeDir).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
