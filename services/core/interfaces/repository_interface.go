package interfaces

import (
	common "altron/common/models"
	"altron/core/models"
	"altron/core/repositories/ent"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*ent.User, error)
	GetUser(ctx context.Context, username string) (*ent.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type ServiceRepository interface {
	CreateService(ctx context.Context, userID uuid.UUID, service *models.Service) (*ent.Service, error)
	GetServiceByID(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID) (*ent.Service, error)
	GetServiceByPort(ctx context.Context, userID uuid.UUID, servicePort uint16) (*ent.Service, error)
	UpdateService(ctx context.Context, userID uuid.UUID, service *models.Service) error
	GetAllServices(ctx context.Context, userID uuid.UUID) ([]*ent.Service, error)
	GetAllPorts(ctx context.Context, userID uuid.UUID) ([]*ent.Service, error)
	DeleteService(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID) error
	GetServicePluginsByPort(ctx context.Context, userID uuid.UUID, servicePort uint16) ([]*ent.Plugin, error)
}

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, userID uuid.UUID, workspace *models.Workspace) (*ent.Workspace, error)
	CreateUserWorkspace(ctx context.Context, workspace *models.Workspace) (*ent.Workspace, error)
	GetWorkspacesByServicePort(ctx context.Context, userID uuid.UUID, servicePort uint16) ([]*ent.Workspace, error)
	GetServiceWorkspaceByName(ctx context.Context, userID uuid.UUID, servicePort uint16, workspaceName string) (*ent.Workspace, error)
	GetWorkspaceByID(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*ent.Workspace, error)
	UpdateWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, workspace *models.Workspace) error
	DeleteWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*ent.Workspace, error)
	UpdateAllWorkspaces(ctx context.Context, userID uuid.UUID, status common.WorkspaceStatus) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, userID uuid.UUID, session *common.Session) error
	CreateSessions(ctx context.Context, workspaceID uuid.UUID, sessions []*common.Session) error
	CreatePcapSession(ctx context.Context, pcapWorkspaceID uuid.UUID, session *common.Session) error
	GetSessionsByWorkspace(ctx context.Context, workspaceID uuid.UUID, filterID *uuid.UUID, paginationIndex int) ([]*ent.Session, error)
	GetPcapSessions(ctx context.Context, pcapWorkspaceID uuid.UUID, paginationIndex int) ([]*ent.Session, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (*ent.Session, error)
	GetSessionsByIDs(ctx context.Context, sessionIDs []uuid.UUID) ([]*ent.Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteAllSessions(ctx context.Context, workspaceID uuid.UUID) error
	CountWorkspaceSessions(ctx context.Context, workspaceID uuid.UUID, filterID *uuid.UUID) (int, error)
	CountSessions(ctx context.Context, userID uuid.UUID) (int, error)
	GetSessions(ctx context.Context, userID uuid.UUID, paginationIndex int) ([]*ent.Session, error)
	GetSessionPackets(ctx context.Context, sessionID uuid.UUID) ([]*ent.Packet, error)
}

type PluginRepository interface {
	CreatePlugins(ctx context.Context, plugin []*common.Plugin) ([]*ent.Plugin, error)
	GetAllPlugins(ctx context.Context) ([]*ent.Plugin, error)
}

type FilterRepository interface {
	CreateFilter(ctx context.Context, userID uuid.UUID, filter *common.Filter) (*ent.Filter, error)
	UpdateFilter(ctx context.Context, userID uuid.UUID, filter *common.Filter) error
	GetAllFilters(ctx context.Context, userID uuid.UUID, servicePort uint16) ([]*ent.Filter, error)
	DeleteFilter(ctx context.Context, userID, filterID uuid.UUID) (*ent.Filter, error)
	CreateSessionFilters(ctx context.Context, sessionFilters []*common.SessionFilter) error
	DeleteSessionFilter(ctx context.Context, userID, filterID uuid.UUID) error
}

type AnalyzerPayloadRepository interface {
	CreateAnalyzerComponents(ctx context.Context, characteristics []string) error
	CreateAnalyzerPayloads(ctx context.Context, workspaceID uuid.UUID, analyzerPayload common.AnalyzerPayload) error
	GetAnalyzerPayloads(ctx context.Context, userID uuid.UUID, servicePort uint16, workspaceID uuid.UUID, componentName string) ([]*ent.AnalyzerPayload, error)
	GetAnalyzerComponents(ctx context.Context) ([]*ent.AnalyzerComponent, error)
}

type CartRepository interface {
	AddSessions(ctx context.Context, workspaceID uuid.UUID, sessionIDs []uuid.UUID) error
	DeleteSession(ctx context.Context, workspaceID uuid.UUID, sessionID uuid.UUID) error
	CountCartSessions(ctx context.Context, workspaceID uuid.UUID) (int, error)
	GetCart(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) ([]*ent.Session, error)
	GetFullCart(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) ([]*ent.Session, error)
	DeleteAllCartSessions(ctx context.Context, workspaceID uuid.UUID) error
}

type PcapWorkspaceRepository interface {
	CreatePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspace *models.PcapWorkspace) (*ent.PcapWorkspace, error)
	UpdatePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID, pcapWorkspace *models.PcapWorkspace) error
	GetAllPcapWorkspaces(ctx context.Context, userID uuid.UUID) ([]*ent.PcapWorkspace, error)
	GetPcapWorkspaceByID(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID) (*ent.PcapWorkspace, error)
	DeletePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID) error
}

type AdminRepository interface {
	ClearAll(ctx context.Context, username string) error
}
