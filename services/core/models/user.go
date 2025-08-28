package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string
	Password string `json:"password"`
}
