package models

import (
	common "altron/common/models"

	"github.com/google/uuid"
)

type PcapWorkspace struct {
	ID       uuid.UUID
	FileName string
	Status   common.WorkspaceStatus
	Sessions []common.Session
}
