package models

import "github.com/google/uuid"

type Plugin struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
