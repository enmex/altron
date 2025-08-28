package interfaces

import "altron/plugin/generated/spec"

type PluginUseCase interface {
	ProcessSession(request *spec.ProcessSessionRequest) (*spec.ProcessSessionResponse, error)
	GetAllPlugins() *spec.GetAllPluginsResponse
}
