package models

import "github.com/google/uuid"

type SessionFilter struct {
	Filter
	SessionID      uuid.UUID `json:"sessionID"`
	MatchesCount   int       `json:"matchesCount"`
	MatchedPackets []int     `json:"matchedPackets"`
}
