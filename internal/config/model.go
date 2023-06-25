package config

type Config struct {
	Commit     CommitConfig
	Overwrites map[string]map[string]string
	Templates  map[string]string
}

type CommitConfig struct {
	Template string
}

func Default() *Config {
	return &Config{
		Commit: CommitConfig{
			Template: "{{.InsertCount}} insertion(s), {{.DeleteCount}} deletion(s)",
		},
		Overwrites: nil,
		Templates:  nil,
	}
}
