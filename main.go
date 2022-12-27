package main

import (
	"log"
	"os"
)

func main() {
	if err := InitApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
