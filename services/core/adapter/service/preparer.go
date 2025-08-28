package service

import (
	common "altron/common/models"
	"altron/core/generated/spec"
	"altron/core/models"
	"github.com/google/uuid"
)

func PrepareUpdateServiceRequest(request *spec.UpdateServiceRequest, serviceID uuid.UUID) *models.Service {
	service := &models.Service{
		ID:          serviceID,
		ContainerID: request.ContainerID,
	}
	if request.Name != nil {
		service.Name = *request.Name
	}
	if request.Link != nil {
		service.Link = *request.Link
	}
	plugins := make([]*common.Plugin, 0, len(request.Plugins))
	for _, plugin := range request.Plugins {
		plugins = append(plugins, &common.Plugin{
			Name: plugin,
		})
	}
	service.Plugins = plugins
	return service
}
