package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                  uuid.UUID                 `json:"id"`
	SentAt              time.Time                 `json:"sentAt"`
	Iface               string                    `json:"iface"`
	ClientHost          string                    `json:"clientHost"`
	ServerPort          uint16                    `json:"serverPort"`
	TTL                 uint8                     `json:"ttl"`
	Protocol            string                    `json:"protocol"`
	Packets             []*Packet                 `json:"packets"`
	PacketsCount        int                       `json:"packetsCount"`
	AverageResponseTime float64                   `json:"averageResponseTime"`
	RequestsNumber      int                       `json:"requestsNumber"`
	MatchedFilters      []*SessionFilter          `json:"matchedFilters"`
	ClientUserAgent     *string                   `json:"clientUserAgent"`
	AnalyzerMatches     map[string]Characteristic `json:"analyzerMatches,omitempty"`
	IsSafe              *bool                     `json:"isSafe,omitempty"`
}
