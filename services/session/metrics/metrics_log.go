package metrics

type MetricsLog struct {
	serversLog map[uint16]int
}

func NewMetricsLog(serversLog map[uint16]int) *MetricsLog {
	return &MetricsLog{
		serversLog: serversLog,
	}
}

func (ml *MetricsLog) ServerRpm(serverPort uint16) int {
	rpm, ok := ml.serversLog[serverPort]
	if !ok {
		return 0
	}
	return rpm
}
