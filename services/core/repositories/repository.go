package repositories

import (
	"altron/core/interfaces"
	"altron/core/repositories/ent"
)

type Repository struct {
	Admin           interfaces.AdminRepository
	User            interfaces.UserRepository
	Service         interfaces.ServiceRepository
	Workspace       interfaces.WorkspaceRepository
	Session         interfaces.SessionRepository
	Plugin          interfaces.PluginRepository
	Filter          interfaces.FilterRepository
	AnalyzerPayload interfaces.AnalyzerPayloadRepository
	Cart            interfaces.CartRepository
	PcapWorkspace   interfaces.PcapWorkspaceRepository
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		Admin:           NewAdminRepository(client),
		User:            NewUserRepository(client),
		Service:         NewServiceRepository(client),
		Workspace:       NewWorkspaceRepository(client),
		Session:         NewSessionRepository(client),
		Plugin:          NewPluginRepository(client),
		Filter:          NewFilterRepository(client),
		AnalyzerPayload: NewAnalyzerPayloadRepository(client),
		Cart:            NewCartRepository(client),
		PcapWorkspace:   NewPcapWorkspaceRepository(client),
	}
}
