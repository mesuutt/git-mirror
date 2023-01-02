package parser

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

type FileStat struct {
	Insert int
	Delete int
	Ext    string
}

type Parser struct {
	allowedTypes map[string]struct{}
}

func NewParser() *Parser {
	return &Parser{}
}

// WithWhitelist sets allowed file types for creating stats from commit diff
// if not given, all file types in the commit will be parsed
func (p *Parser) WithWhitelist(whitelist []string) *Parser {
	if whitelist == nil {
		return p
	}

	// converting to map for find with O(1)
	types := make(map[string]struct{}, len(whitelist))
	for _, v := range whitelist {
		types[strings.TrimSpace(v)] = struct{}{}
	}

	p.allowedTypes = types

	return p
}

// Parse parses diff output and create stat for each file type
func (p *Parser) Parse(r io.Reader) ([]FileStat, error) {
	groupedStats := make(map[string]FileStat)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		ext := filepath.Ext(parts[len(parts)-1])
		var addCount, delCount int
		var err error

		if addCount, err = strconv.Atoi(parts[0]); err != nil {
			return nil, fmt.Errorf("add count parse failed: `%s`", line)
		}

		if delCount, err = strconv.Atoi(parts[1]); err != nil {
			return nil, fmt.Errorf("delete count parse failed: `%s`", line)
		}

		if ext == "" {
			// If file does not have an extension
			// search in allowed types and use filename
			// if allowed types not given that means we can use filename too
			// eg: Makefile
			filename := parts[len(parts)-1]
			if p.allowedTypes == nil {
				ext = filename
			} else if _, ok := p.allowedTypes[filename]; ok {
				ext = filename
			}
		}

		if p.allowedTypes != nil {
			if _, ok := p.allowedTypes[strings.Replace(ext, ".", "", 1)]; !ok {
				continue
			}
		}

		s, ok := groupedStats[ext]
		if !ok {
			groupedStats[ext] = FileStat{
				Insert: addCount,
				Delete: delCount,
				Ext:    ext,
			}
			continue
		}

		s.Insert += addCount
		s.Delete += delCount
		groupedStats[ext] = s
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var stats []FileStat
	for i := range groupedStats {
		stats = append(stats, groupedStats[i])
	}

	return stats, nil
}
