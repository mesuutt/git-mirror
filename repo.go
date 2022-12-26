package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Repo struct {
	path string
}

func NewRepo(path string) Repo {
	return Repo{path: path}
}

func (r Repo) AddStat(stat FileStat) error {
	now := time.Now()
	dir := filepath.Join(r.path, strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()))

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("directory could not created: `%s`. error: %w", dir, err)
	}

	filePath := filepath.Join(dir, fmt.Sprintf("log%s", stat.Ext))

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("file could not be opened for appending changes: `%s`. error: %w", filePath, err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%d insertion(s), %d deletion(s)", stat.Insert, stat.Delete))
	if err != nil {
		return fmt.Errorf("changes could not write to file: `%s`. error: %w", filePath, err)
	}

	return nil
}
