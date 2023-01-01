package gitmirror

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	input := `
1       2       a.go
1       0       a.txt
1       3       c.rs
1       0       foo bar baz.rs
19      6       bar.json
10      2       app.py
1       2       b.go
1       2       Makefile
`
	p := NewParser()
	result, err := p.Parse(strings.NewReader(input))

	if err != nil {
		t.Fatal(err)
	}

	expected := map[string][]int{
		".go":      {2, 4},
		".rs":      {2, 3},
		".txt":     {1, 0},
		".json":    {19, 6},
		".py":      {10, 2},
		"Makefile": {1, 2},
	}

	for _, v := range result {
		if v.Insert != expected[v.Ext][0] {
			t.Errorf("additon not mathes for %s: want: %d, got: %d", v.Ext, expected[v.Ext][0], v.Insert)
		}

		if v.Delete != expected[v.Ext][1] {
			t.Errorf("deletion not mathes for %s: want: %d, got: %d", v.Ext, expected[v.Ext][1], v.Delete)
		}
	}
}

func TestParseWithAllowedFileTypes(t *testing.T) {
	input := `
1       2       a.go
1       0       a.txt
1       2       Makefile
`
	tests := []struct {
		name            string
		allowedTypes    []string
		expectedStatLen int
	}{
		{
			name:            "types given",
			allowedTypes:    []string{"go", "Makefile"},
			expectedStatLen: 2,
		},
		{
			name:            "types not given",
			allowedTypes:    nil,
			expectedStatLen: 3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			p := NewParser().WithWhitelist(tt.allowedTypes)
			result, err := p.Parse(strings.NewReader(input))
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != tt.expectedStatLen {
				t.Errorf("parsed file change count is different. want %d, got: %d", tt.expectedStatLen, len(result))
			}
		})
	}
}
