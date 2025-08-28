package config

import (
	"altron/pkg/plugin"
	"os"
)

func NewPluginConfig() *plugin.Config {
	return &plugin.Config{
		PluginsDirectory: os.Getenv("PLUGINS_DIRECTORY"),
	}
}
