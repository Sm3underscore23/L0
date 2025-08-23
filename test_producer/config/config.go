package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	// Config -.
	Config struct {
		Kafka Kafka `yaml:"kafka"`
	}

	Kafka struct {
		BrokerList []string `yaml:"broker_list"`
		Topic      string   `yaml:"topic"`
	}
)

// NewConfig returns app config.
func NewConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}

	rowConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(rowConfig, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
