package repositories

import (
	"altron/common/models"
	"altron/core/interfaces"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/plugin"
	"context"
	"strings"
)

var _ interfaces.PluginRepository = (*PluginRepository)(nil)

type PluginRepository struct {
	Client *ent.Client
}

func NewPluginRepository(client *ent.Client) *PluginRepository {
	return &PluginRepository{
		Client: client,
	}
}

func (r *PluginRepository) CreatePlugins(ctx context.Context, plugins []*models.Plugin) ([]*ent.Plugin, error) {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	pluginsEnt := make([]*ent.Plugin, 0, len(plugins))

	for _, p := range plugins {
		_, err := tx.Plugin.
			Query().
			Where(plugin.NameEQ(p.Name)).
			Only(ctx)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				if err := tx.Rollback(); err != nil {
					return nil, err
				}
				return nil, err
			}
			pluginEnt, err := tx.Plugin.
				Create().
				SetName(p.Name).
				Save(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return nil, err
				}
				return nil, err
			}
			pluginsEnt = append(pluginsEnt, pluginEnt)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pluginsEnt, nil
}

func (r *PluginRepository) GetAllPlugins(ctx context.Context) ([]*ent.Plugin, error) {
	plugins, err := r.Client.Plugin.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}
