package models

import "time"

type Packet struct {
	SentAt    time.Time `json:"sentAt"`
	Payload   string    `json:"payload"`
	IsRequest bool      `json:"isRequest"`
}
