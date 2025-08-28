package interfaces

type PluginInterface interface {
	Process(sessionPayload [][]byte) ([][]byte, error)
}
