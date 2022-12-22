package main

import (
	"encoding/json"
	"log"
	"strings"
)

func main() {
	/*out, err := exec.Command("git", " diff", "--stat", "HEAD~1").Output()
	if err != nil {
		log.Fatal(err)
	}*/

	input := `21122022.py | 1 +
 21122022.rs | 2 ++
 multi sss.txt | 1 +
 21122022.go   | 5 ++++-
 21122022.txt   | 1 +
 21122022.json   | 2 ++
 21122021.go   | 3 ++-
 multi sss.properties | 1 +`

	result, err := ParseDiff(strings.NewReader(input))
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(result)

	println("sss", string(b))
}
