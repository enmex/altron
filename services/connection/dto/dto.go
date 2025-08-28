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

type ClientDto struct {
	Action  string             `json:"action"`
	Payload utils.Pair[string] `json:"payload,omitempty"`
}

type GetCheckerMaskResponse struct {
	HasCheckerMask bool                               `json:"hasCheckerMask"`
	CheckerMask    map[string][]common.Characteristic `json:"checkerMask,omitempty"`
}

type AddPcapSessionRequest struct {
	Session *common.Session `json:"session"`
}

type UpdatePcapWorkspaceStatusRequest struct {
	Status common.WorkspaceStatus `json:"status"`
}

type ServiceResponse struct {
	Port        uint16  `json:"port"`
	ContainerID *string `json:"containerID,omitempty"`
}

type AgentInfoResponse struct {
	ContainerID string `json:"containerID"`
}

type GetAllServicesResponse struct {
	Services []ServiceResponse `json:"services"`
}
