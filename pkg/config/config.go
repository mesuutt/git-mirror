package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Commit     CommitConfig
	Overwrites map[string]map[string]string
	Templates  map[string]string
}

type CommitConfig struct {
	Template string
}

func ReadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := toml.Unmarshal(bytes, &conf); err != nil {
		return nil, fmt.Errorf("config file parse failed: %s, error: %w", path, err)
	}

	return &conf, nil
}
