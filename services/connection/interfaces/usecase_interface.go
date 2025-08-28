package interfaces

import (
	"altron/connection/generated/spec"
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionUseCase interface {
	ListenSessions(ctx context.Context, ws *websocket.Conn, servicePort uint16) error
	ListenLogs(ctx context.Context, ws *websocket.Conn, containerID string) error
	ListenStats(ctx context.Context, ws *websocket.Conn) error
	OpenMainChannel(ctx context.Context, ws *websocket.Conn) error
	GetAgentStatus(ctx context.Context) *spec.GetAgentStatusResponse
	CreateEvent(ctx context.Context, request *spec.CreateEventRequest) (*spec.CreateEventResponse, error)
}

type WorkspaceUseCase interface {
	StartListeningSessions(ctx context.Context, workspaceID uuid.UUID, servicePort uint16, request *spec.StartSessionsListeningRequest) error
	StopListeningSessions(workspaceID uuid.UUID)
}

type PcapWorkspaceUseCase interface {
	StartListeningPcap(pcapWorkspaceID uuid.UUID, request *spec.StartPcapListeningRequest)
}
