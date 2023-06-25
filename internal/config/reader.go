package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

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
