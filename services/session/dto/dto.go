package dto

import (
	common "altron/common/models"
	"altron/utils"
)

type AddSessionsToWorkspaceRequest struct {
	Sessions []common.Session
}

type ProcessSessionRequest struct {
	Session *common.Session
	Plugins []string
}

type ProcessSessionResponse = common.Session

type ApplyFiltersRequest struct {
	Filters []*common.Filter
}

type GetAllFiltersResponse struct {
	Filters []*common.Filter `json:"filters"`
}

type GetServicePluginsResponse struct {
	Plugins []string `json:"plugins"`
}

type UpdateWorkspaceStatusRequest struct {
	Status common.WorkspaceStatus `json:"status"`
}

type CreateAnalyzerPayloadRequest struct {
	Payload map[string][]common.Characteristic `json:"payload"`
}

type ClientRequest struct {
	Action  string             `json:"action"`
	Payload utils.Pair[string] `json:"payload,omitempty"`
}

type SendServerMetricsResponse struct {
	Rpm int `json:"rpm"`
}

type GetUserInfoResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServiceResponse struct {
	Port        uint16  `json:"port"`
	ContainerID *string `json:"containerID,omitempty"`
}

type GetPortsResponse struct {
	Services []ServiceResponse `json:"services"`
}
