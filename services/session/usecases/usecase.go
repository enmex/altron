package usecases

import "altron/session/interfaces"

type UseCase struct {
	Packet           interfaces.PacketUseCase
	SessionCollector interfaces.SessionCollectorUseCase
	LogsCollector    interfaces.LogsCollectorUseCase
}

func NewUseCase(
	packet interfaces.PacketUseCase,
	sessionTree interfaces.SessionCollectorUseCase,
	logsCollector interfaces.LogsCollectorUseCase,
) *UseCase {
	return &UseCase{
		Packet:           packet,
		SessionCollector: sessionTree,
		LogsCollector:    logsCollector,
	}
}
