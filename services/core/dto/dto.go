package dto

import "altron/core/models"

type PortRequest struct {
	Port        uint16  `json:"port"`
	ContainerID *string `json:"containerID,omitempty"`
}

type CreatePortsRequest struct {
	Ports []PortRequest `json:"ports"`
}

type DeletePortRequest struct {
	PortRequest
}

type GetPortsResponse struct {
	Ports []PortRequest `json:"ports"`
}

type UpdateContainerRequest struct {
	OldContainer *string `json:"oldContainer,omitempty"`
	NewContainer string  `json:"newContainer"`
}

type StartSessionsListeningsRequest struct {
	Timeout   string `json:"timeout"`
	StartTime *int64 `json:"startTime,omitempty"`
	IsChecker bool   `json:"isChecker,omitempty"`
}

type ConvertPacketToExploitRequest struct {
	Base64Str string `json:"base64_str"`
}

type ConvertToExploitResponse struct {
	Result string `json:"result"`
}

type ConvertSessionToExploitRequest struct {
	Base64Strings []string `json:"base64_strings"`
}

type ExtractFilesRequest struct {
	Base64Str    string `json:"base64_str"`
	SessionID    string `json:"session_id"`
	PacketNumber string `json:"packet_number"`
}

type StartPcapListeningRequest struct {
	FileName string `json:"fileName"`
}

type UploadPcapRequest struct {
	FileName string `json:"fileName"`
}

type CreateEventRequest = models.Event

type CreateEventResponse struct {
	Status string `json:"status"`
	Data   []byte `json:"data,omitempty"`
}

type ScanServiceRequest struct {
	Port uint16 `json:"port"`
}

type ScanHostServiceRequest struct {
	Port uint16 `json:"port"`
}

type ScanHostServiceResponse struct {
	Link string `json:"link"`
}

type AddPortRequest struct {
	Port        uint16  `json:"port"`
	ContainerID *string `json:"containerID,omitempty"`
}

type SFTPTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ChangeUserPasswordRequest struct {
	HomeDir     string              `json:"home_dir"`
	Password    string              `json:"password"`
	Permissions map[string][]string `json:"permissions"`
	Status      int                 `json:"status"`
}
