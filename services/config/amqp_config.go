package config

import (
	"altron/pkg/amqp"
	"os"

	"gopkg.in/yaml.v2"
)

func NewAMQPConfig() (*amqp.Config, error) {
	cfg := amqp.Config{}

	yamlFile, err := os.ReadFile("/config.yaml")
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}
	cfg.Prefix = os.Getenv("HOSTNAME")
	cfg.Host = os.Getenv("RABBITMQ_SERVICE_HOST")

	return &cfg, nil
}
