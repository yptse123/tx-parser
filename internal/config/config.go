package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Port   string `yaml:"port"`
	Host   string `yaml:"host"`
	Ethrpc string `yaml:"ethrpc"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
