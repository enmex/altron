package models

type EventType string

const (
	ScanHost       EventType = "scan_localhost"
	ScanContainers EventType = "scan_containers"
	ScanPort       EventType = "scan_port"
	CreateService  EventType = "create_service"
	DeleteService  EventType = "delete_service"
	UpdateService  EventType = "update_service"
	GetContainers  EventType = "get_containers"
	AltronReset    EventType = "altron_reset"
)

type Event struct {
	Type            EventType `json:"type"`
	Data            []byte    `json:"data,omitempty"`
	WaitForResponse bool      `json:"waitForResponse,omitempty"`
}
