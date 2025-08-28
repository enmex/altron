package usecases

import (
	common "altron/common/models"
	"altron/core/adapter/plugin"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"context"
)

var _ interfaces.PluginUseCase = (*PluginUseCase)(nil)

type PluginUseCase struct {
	repo interfaces.PluginRepository
}

func NewPluginUseCase(repo interfaces.PluginRepository) *PluginUseCase {
	return &PluginUseCase{
		repo: repo,
	}
}

func (u *PluginUseCase) CreatePlugins(ctx context.Context, plugins []string) error {
	pluginsModel := make([]*common.Plugin, 0, len(plugins))
	for _, plugin := range plugins {
		pluginsModel = append(pluginsModel, &common.Plugin{
			Name: plugin,
		})
	}

	_, err := u.repo.CreatePlugins(ctx, pluginsModel)
	return err
}

func (u *PluginUseCase) GetAllPlugins(ctx context.Context) (*spec.GetAllPluginsResponse, error) {
	pluginsEnt, err := u.repo.GetAllPlugins(ctx)
	if err != nil {
		return nil, err
	}
	return &spec.GetAllPluginsResponse{
		Plugins: plugin.PresentPlugins(pluginsEnt),
	}, nil
}
