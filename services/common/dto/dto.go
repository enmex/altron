package dto

import "altron/common/models"

type SendPacketToQueueRequest struct {
	FileName   *string
	Interface  string
	ServerPort uint16
	ClientHost string
	Protocol   string
	TTL        uint8
	Packet     models.Packet
	IsLast     bool
}

const (
	ErrorBadRequest          = "BAD_REQUEST"
	ErrorInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrorUnauthorized        = "UNAUTHORIZED"
	ErrorForbidden           = "FORBIDDEN"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SelectedCharacteristicDto struct {
	Value    string `json:"value"`
	Selected bool   `json:"selected"`
	Blocked  bool   `json:"blocked"`
}

type GetServiceCheckerMaskResponse struct {
	Analyzer map[string][]models.Characteristic `json:"analyzer"`
}
