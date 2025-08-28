package config

import (
	"net"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AltronName           string
	ManagerPassword      string        `yaml:"manager-password"`
	TcpStreamTimeout     time.Duration `yaml:"tcp-stream-timeout"`
	UdpStreamTimeout     time.Duration `yaml:"udp-stream-timeout"`
	AltronHost           string
	AltronPort           int               `yaml:"altron-port"`
	AltronPluginPort     int               `yaml:"altron-plugin-port"`
	AltronSessionPort    int               `yaml:"altron-session-port"`
	AltronConverterPort  int               `yaml:"altron-converter-port"`
	AltronConnectionPort int               `yaml:"altron-connection-port"`
	AltronContainers     map[string]string `yaml:"altron-containers"`
}

func NewAppConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	yamlFile, err := os.ReadFile("/config.yaml")
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}
	cfg.AltronName = os.Getenv("HOSTNAME")
	cfg.ManagerPassword = os.Getenv("MANAGER_PASSWORD")

	addresses, err := net.LookupHost(os.Getenv("HOSTNAME"))
	if err != nil {
		return nil, err
	}
	cfg.AltronHost = addresses[0]

	return cfg, nil
}
