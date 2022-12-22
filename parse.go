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
	Change int
	Ext    string
}

func ParseDiff(r io.Reader) ([]FileStat, error) {
	groupedStats := make(map[string]FileStat)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		ext := filepath.Ext(parts[len(parts)-4])
		changeCount, err := strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			return nil, fmt.Errorf("change count parse failed: `%s`", line)
		}

		s, ok := groupedStats[ext]
		if !ok {
			groupedStats[ext] = FileStat{
				Change: changeCount,
				Ext:    ext,
			}
			continue
		}

		s.Change += changeCount
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
