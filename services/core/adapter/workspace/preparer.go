package workspace

import (
	common "altron/common/models"
	"altron/core/generated/spec"
	"altron/core/models"
)

func PrepareCreateWorkspaceRequest(request *spec.CreateWorkspaceRequest) *models.Workspace {
	var status common.WorkspaceStatus
	if request.StartTime != nil {
		status = common.WAITING
	} else {
		status = common.LISTENING
	}
	return &models.Workspace{
		Name:        request.Name,
		ServicePort: request.ServicePort,
		Status:      status,
	}
}

func PrepareUpdateWorkspaceRequest(request *spec.UpdateWorkspaceRequest) *models.Workspace {
	var status common.WorkspaceStatus
	if request.Status != nil {
		status = common.WorkspaceStatus(*request.Status)
	}
	var name string
	if request.Name != nil {
		name = *request.Name
	}
	return &models.Workspace{
		Name:   name,
		Status: status,
	}
}
