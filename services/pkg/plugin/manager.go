package plugin

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"sort"
	"strings"
)

type PluginManager[T any] struct {
	Plugins   map[int]T
	loadOrder map[string]int
}

func NewPluginManager[T any](cfg *Config) (*PluginManager[T], error) {
	files, err := ioutil.ReadDir(cfg.PluginsDirectory)
	if err != nil {
		return nil, err
	}
	pluginsLoadOrderBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/plugins_order.txt", cfg.PluginsDirectory))
	if err != nil {
		return nil, err
	}

	pluginsLoadOrder := make(map[string]int, 0)
	order := strings.Split(string(pluginsLoadOrderBytes), "\n")
	if len(order) != 0 {
		for idx, plugin := range order {
			pluginsLoadOrder[plugin] = idx
		}
	}
	plugins := make(map[int]T)

	for _, file := range files {
		if file.IsDir() {
			p, err := plugin.Open(fmt.Sprintf("%s/%s/%s.so", cfg.PluginsDirectory, file.Name(), file.Name()))
			if err != nil {
				return nil, err
			}
			s, err := p.Lookup("Main")
			if err != nil {
				return nil, err
			}
			plugin, ok := s.(T)
			if !ok {
				return nil, fmt.Errorf("plugin %s is not suitable", file.Name())
			}
			plugins[pluginsLoadOrder[file.Name()]] = plugin
		}
	}

	return &PluginManager[T]{
		Plugins:   plugins,
		loadOrder: pluginsLoadOrder,
	}, nil
}

func (u *PluginManager[T]) GetPluginsByNames(names []string) []T {
	orders := make([]int, 0, len(names))
	plugins := make([]T, 0, len(names))

	for _, name := range names {
		orders = append(orders, u.loadOrder[name])
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i] < orders[j]
	})
	for _, order := range orders {
		plugins = append(plugins, u.Plugins[order])
	}
	return plugins
}

func (u *PluginManager[T]) GetPluginNames() []string {
	names := make([]string, 0)
	for name := range u.loadOrder {
		names = append(names, name)
	}
	return names
}
