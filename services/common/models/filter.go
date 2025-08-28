package models

import "github.com/google/uuid"

type Filter struct {
	ID           uuid.UUID  `json:"id"`
	ServiceID    *uuid.UUID `json:"serviceId,omitempty"`
	Name         string     `json:"name"`
	Color        string     `json:"color"`
	Regex        *string    `json:"regex,omitempty"`
	TTL          *uint8     `json:"ttl,omitempty"`
	TotalPackets *int       `json:"totalPackets,omitempty"`
	IsBlocking   bool       `json:"isBlocking"`
	InRequest    bool       `json:"inRequest"`
	InResponse   bool       `json:"inResponse"`
}
