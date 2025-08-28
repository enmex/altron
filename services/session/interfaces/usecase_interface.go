package interfaces

import (
	"altron/session/generated/spec"
	"context"

	"github.com/google/gopacket"
)

type PacketUseCase interface {
	ProducePacket(ctx context.Context, packet gopacket.Packet, fileName *string)
	AddPcapPorts(ctx context.Context, request *spec.CreatePortsRequest) error
	DeletePcapPort(ctx context.Context, request *spec.DeletePortRequest) error
	ImportFile(ctx context.Context, request *spec.UploadPcapRequest) error
}

type SessionCollectorUseCase interface {
	DeleteServer(serverPort uint16) error
	GetServers() []uint16
	ServerExists(serverPort uint16) bool
	AddServer(serverPort uint16) error
}

type LogsCollectorUseCase interface {
	DeleteLogsContainer(containerID string) error
	AddLogsContainer(containerID string) error
	UpdateContainer(ctx context.Context, request *spec.UpdateContainerRequest) error
	GetCachedLogs(containerID string) (*spec.GetLatestContainerLogsResponse, error)
}

type HealthUseCase interface {
	LaunchHealthMonitoring(ctx context.Context) error
}
