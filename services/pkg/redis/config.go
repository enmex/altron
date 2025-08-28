package redis

import "time"

type Config struct {
	Host           string        `yaml:"redis-host"`
	Port           int           `yaml:"redis-port"`
	User           string        `yaml:"redis-user"`
	Password       string        `yaml:"redis-password"`
	ExpirationTime time.Duration `yaml:"redis-expiration-time"`
	Prefix         string
}
