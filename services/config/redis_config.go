package config

import (
	"altron/pkg/redis"
	"os"

	"gopkg.in/yaml.v2"
)

func NewRedisConfig() (*redis.Config, error) {
	cfg := redis.Config{}

	yamlFile, err := os.ReadFile("/config.yaml")
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}
	cfg.Prefix = os.Getenv("HOSTNAME")
	cfg.Host = os.Getenv("REDIS_MASTER_SERVICE_HOST")

	return &cfg, nil
}
