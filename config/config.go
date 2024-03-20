package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

type Config struct {
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	yamlfile, _ := os.ReadFile(configPath)
	if err := yaml.Unmarshal(yamlfile, &cfg); err != nil {
		log.Fatalf("config file is incorrect: %v", err)
	}
	return &cfg
}
