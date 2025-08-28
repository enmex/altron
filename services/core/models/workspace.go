package models

import (
	common "altron/common/models"

	"github.com/google/uuid"
)

type Workspace struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	ServicePort uint16                 `json:"service_port"`
	Sessions    []common.Session       `json:"sessions"`
	Status      common.WorkspaceStatus `json:"status"`
}
