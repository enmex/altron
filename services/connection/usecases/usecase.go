package usecases

import "altron/connection/interfaces"

type UseCase struct {
	Connection    interfaces.ConnectionUseCase
	Workspace     interfaces.WorkspaceUseCase
	PcapWorkspace interfaces.PcapWorkspaceUseCase
}

func NewUseCase(
	connection interfaces.ConnectionUseCase,
	workspace interfaces.WorkspaceUseCase,
	pcapWorkspace interfaces.PcapWorkspaceUseCase,
) *UseCase {
	return &UseCase{
		Connection:    connection,
		Workspace:     workspace,
		PcapWorkspace: pcapWorkspace,
	}
}
