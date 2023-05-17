package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Commit    commitConfig
	Aliases   map[string]map[string]string
	Templates map[string]string
}

type commitConfig struct {
	Template string
}

func ReadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := toml.Unmarshal(bytes, &conf); err != nil {
		fmt.Printf("config file parse failed: %s, error: %v\n", path, err)
		return nil, err
	}

	return &conf, nil
}
