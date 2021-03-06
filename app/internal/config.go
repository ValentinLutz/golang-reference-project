package internal

import (
	"app/external/database"
	"app/internal/config"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Region      config.Region      `yaml:"region"`
	Environment config.Environment `yaml:"environment"`
	Server      ServerConfig       `yaml:"server"`
	Logger      LoggerConfig       `yaml:"logger"`
	Database    database.Config    `yaml:"database"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout TimeoutConfig `yaml:"timeout"`
}

type TimeoutConfig struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

func NewConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	var decodedConfig Config
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		return Config{}, err
	}

	return decodedConfig, nil
}
