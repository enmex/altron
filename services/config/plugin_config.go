package config

import (
	"altron/pkg/plugin"
)

func NewPluginConfig() *plugin.Config {
	return &plugin.Config{
		PluginsDirectory: "./plugins", //os.Getenv("PLUGINS_DIRECTORY"),
	}
}
