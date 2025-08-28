package http

import (
	"net"
)

const DefaultPort = "8080"

type Config struct {
	Port string `json:"port"`
}

// Addr returns server address in format ":<port>"
func (c *Config) Addr() string {
	return net.JoinHostPort("", c.Port)
}
