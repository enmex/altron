package config

import (
	"altron/pkg/db"
	"os"

	"gopkg.in/yaml.v2"
)

func NewDBConfig() (*db.Config, error) {
	cfg := db.Config{}

	yamlFile, err := os.ReadFile("/config.yaml")
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}
	cfg.Host = os.Getenv("POSTGRESQL_SERVICE_HOST")
	return &cfg, nil
}
