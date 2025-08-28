package plugin

import "altron/core/repositories/ent"

func PresentPlugins(pluginsEnt []*ent.Plugin) []string {
	plugins := make([]string, 0, len(pluginsEnt))
	for _, pluginEnt := range pluginsEnt {
		plugins = append(plugins, pluginEnt.Name)
	}
	return plugins
}
