package main

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

func NewParser(allowedTypes []string) Parser {
	// converting to map for find with O(1)
	types := make(map[string]struct{}, len(allowedTypes))
	for _, v := range allowedTypes {
		types[v] = struct{}{}
	}

	return Parser{allowedTypes: types}
}

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
		addCount, err := strconv.Atoi(parts[0])
		delCount, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("stat parse failed: `%s`", line)
		}

		if ext == "" {
			// If file does not have an extension
			// search in allowed types and use filename
			// eg: Makefile
			filename := parts[len(parts)-1]
			if _, ok := p.allowedTypes[filename]; ok {
				ext = filename
			}
		}

		if _, ok := p.allowedTypes[ext]; !ok {
			continue
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
