package usecases

import "altron/core/interfaces"

type UseCase struct {
	Admin           interfaces.AdminUseCase
	Auth            interfaces.AuthUseCase
	Dashboard       interfaces.DashboardUseCase
	Service         interfaces.ServiceUseCase
	Workspace       interfaces.WorkspaceUseCase
	Session         interfaces.SessionUseCase
	Conversion      interfaces.ConversionUseCase
	Plugin          interfaces.PluginUseCase
	Filter          interfaces.FilterUseCase
	SessionAnalyzer interfaces.SessionAnalyzerUseCase
	Cart            interfaces.CartUseCase
	PcapWorkspace   interfaces.PcapWorkspaceUseCase
}

func NewUseCase(
	admin interfaces.AdminUseCase,
	auth interfaces.AuthUseCase,
	dashboard interfaces.DashboardUseCase,
	service interfaces.ServiceUseCase,
	workspace interfaces.WorkspaceUseCase,
	session interfaces.SessionUseCase,
	conversion interfaces.ConversionUseCase,
	plugin interfaces.PluginUseCase,
	filter interfaces.FilterUseCase,
	sessionAnalyzer interfaces.SessionAnalyzerUseCase,
	cart interfaces.CartUseCase,
	pcapWorkspace interfaces.PcapWorkspaceUseCase,
) *UseCase {
	return &UseCase{
		Admin:           admin,
		Auth:            auth,
		Dashboard:       dashboard,
		Service:         service,
		Workspace:       workspace,
		Session:         session,
		Conversion:      conversion,
		Plugin:          plugin,
		Filter:          filter,
		SessionAnalyzer: sessionAnalyzer,
		Cart:            cart,
		PcapWorkspace:   pcapWorkspace,
	}
}
