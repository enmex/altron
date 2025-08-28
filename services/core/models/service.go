package models

import (
	"altron/common/models"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Link        string           `json:"link"`
	Port        uint16           `json:"port"`
	ContainerID *string          `json:"containerID,omitempty"`
	Plugins     []*models.Plugin `json:"plugins"`
}
