package main

import (
	"strings"
	"testing"
)

func TestParseDiff(t *testing.T) {
	input := `
app.py | 1 +
 magic.rs | 2 ++
 foo.go   | 5 ++++-
 bar baz.txt   | 1 +
 my config.json   | 2 ++
 other go file.go   | 3 ++-
 my journey.txt | 1 +
`
	result, err := ParseDiff(strings.NewReader(input))

	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]int{
		".go":   8,
		".rs":   2,
		".txt":  2,
		".json": 2,
		".py":   1,
	}

	for _, v := range result {
		if v.Change != expected[v.Ext] {
			t.Errorf("diff count not mathes for %s: want: %d, got: %d", v.Ext, expected[v.Ext], v.Change)
		}
	}
}
