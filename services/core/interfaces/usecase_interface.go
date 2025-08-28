package interfaces

import (
	"altron/core/dto"
	"altron/core/generated/spec"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, request *spec.AuthRequest) (*spec.AuthResponse, error)
	SignIn(ctx context.Context, request *spec.AuthRequest) (*spec.AuthResponse, error)
	SignUpManager(ctx context.Context) (*spec.AuthResponse, error)
	SignInManager(ctx context.Context, request *spec.AuthRequest) (*spec.AuthResponse, error)
	AuthUser(ctx context.Context) (*spec.AuthResponse, error)
	GetUserInfo(ctx context.Context) (*spec.GetUserInfoResponse, error)
}

type DashboardUseCase interface {
	GetDashboard(ctx context.Context, userID uuid.UUID) (*spec.GetDashboardResponse, error)
}

type ServiceUseCase interface {
	CreateService(ctx context.Context, userID uuid.UUID, request *spec.CreateServiceRequest) (*spec.CreateServiceResponse, error)
	GetService(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetServiceResponse, error)
	UpdateService(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID, request *spec.UpdateServiceRequest) error
	GetAllPorts(ctx context.Context, userID uuid.UUID) (*dto.GetPortsResponse, error)
	DeleteService(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID) error
	ScanHost(ctx context.Context, scope string) (*spec.ScanHostServicesResponse, error)
	ScanService(ctx context.Context, port uint16) (*spec.ScanHostServiceResponse, error)
}

type WorkspaceUseCase interface {
	CreateWorkspace(ctx context.Context, userID uuid.UUID, request *spec.CreateWorkspaceRequest) (*spec.CreateWorkspaceResponse, error)
	GetWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*spec.GetWorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, request *spec.UpdateWorkspaceRequest) error
	DeleteWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) error
	ResetWorkspaces(ctx context.Context, userID uuid.UUID) error
}

type SessionUseCase interface {
	AddSessions(ctx context.Context, workspaceID uuid.UUID, request *spec.AddSessionsToWorkspaceRequest) error
	AddPcapSession(ctx context.Context, pcapWorkspace uuid.UUID, request *spec.AddPcapSessionRequest) error
	GetPaginatedSessions(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) (*spec.GetSessionsResponse, error)
	GetPaginatedPcapSessions(ctx context.Context, pcapWorkspaceID uuid.UUID, paginationIndex int) (*spec.GetSessionsResponse, error)
	SearchSessions(ctx context.Context, workspaceID uuid.UUID, paginationIndex int, request *spec.SearchSessionsRequest) (*spec.SearchSessionsResponse, error)
	CreateSession(ctx context.Context, userID uuid.UUID, request *spec.CreateSessionRequest) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*spec.GetSessionResponse, error)
	ClearSessions(ctx context.Context, workspaceID uuid.UUID) error
	MergeSessions(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*spec.MergeSessionsResponse, error)
}

type CartUseCase interface {
	AddSessions(ctx context.Context, workspaceID uuid.UUID, request *spec.AddSessionsToCartRequest) error
	GetCart(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) (*spec.GetCartResponse, error)
	RemoveSession(ctx context.Context, workspaceID, sessionID uuid.UUID) error
	ClearCart(ctx context.Context, workspaceID uuid.UUID) error
}

type ConversionUseCase interface {
	ConvertSessionToExploit(ctx context.Context, request *spec.ConvertSessionToExploitRequest) (*spec.ConvertToExploitResponse, error)
	ConvertPacketToExploit(ctx context.Context, request *spec.ConvertPacketToExploitRequest) (*spec.ConvertToExploitResponse, error)
	ExtractFiles(ctx context.Context, request *spec.ExtractFilesFromPacketRequest) (*spec.ExtractFilesResponse, error)
}

type PluginUseCase interface {
	CreatePlugins(ctx context.Context, plugins []string) error
	GetAllPlugins(ctx context.Context) (*spec.GetAllPluginsResponse, error)
}

type FilterUseCase interface {
	CreateFilter(ctx context.Context, userID uuid.UUID, request *spec.CreateFilterRequest) error
	UpdateFilter(ctx context.Context, userID uuid.UUID, request *spec.UpdateFilterRequest, filterID uuid.UUID) error
	GetAllFilters(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetAllFiltersResponse, error)
	DeleteFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID) error
}

type SessionAnalyzerUseCase interface {
	GetAllAnalyzerComponents(ctx context.Context) (*spec.GetAnalyzerComponentsResponse, error)
	CreateAnalyzerPayload(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, request *spec.CreateAnalyzerPayloadRequest) error
	GetAnalyzerPayload(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetSessionAnalyzerResponse, error)
	GetWorkspaceAnalyzerPayload(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*spec.GetSessionAnalyzerResponse, error)
	GetServiceCheckerMask(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetSessionAnalyzerResponse, error)
}

type PcapWorkspaceUseCase interface {
	CreatePcapWorkspace(ctx context.Context, userID uuid.UUID, file multipart.File, fileHeader *multipart.FileHeader) (*spec.CreatePcapWorkspaceResponse, error)
	UpdatePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID, request *spec.UpdatePcapWorkspaceRequest) error
	GetPcapWorkspace(ctx context.Context, userID, pcapWorkspaceID uuid.UUID) (*spec.GetPcapWorkspaceResponse, error)
	DeletePcapWorkspace(ctx context.Context, userID, pcapWorkspaceID uuid.UUID) error
}

type AdminUseCase interface {
	ResetAll(ctx context.Context) (*spec.ReconfigureAltronResponse, error)
	GetInfo(ctx context.Context) (*spec.GetInfoResponse, error)
}
