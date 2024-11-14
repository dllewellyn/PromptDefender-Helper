package gh

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Prompts []string `yaml:"prompts"`
}

func LoadConfigFromString(content string) (*Config, error) {
	var config Config
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config content: %w", err)
	}

	return &config, nil
}
