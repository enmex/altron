package metrics

import (
	"altron/pkg/amqp"
	"altron/session/dto"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type ServerMetrics struct {
	log                  *logrus.Logger
	serverPortsStatistic sync.Map //*utils.ConcurrentMap[uint16, int] //map[uint16]int
	lastLog              *MetricsLog
	producer             *amqp.Producer
}

func NewServerMetrics(log *logrus.Logger, client *amqp.Client) (*ServerMetrics, error) {
	producer, err := amqp.NewProducer(client)
	if err != nil {
		return nil, err
	}

	return &ServerMetrics{
		log:                  log,
		serverPortsStatistic: sync.Map{},
		producer:             producer,
	}, nil
}

func (sm *ServerMetrics) Checkpoint() {
	metricsLog := make(map[uint16]int)
	sm.serverPortsStatistic.Range(func(key, value any) bool {
		metricsLog[key.(uint16)] = value.(int)
		return true
	})

	sm.lastLog = NewMetricsLog(metricsLog)
	for serverPort := range metricsLog {
		sm.serverPortsStatistic.LoadOrStore(serverPort, 0)
	}
}

func (sm *ServerMetrics) SendMetrics() error {
	if sm.lastLog == nil {
		return nil
	}
	lastMetrics := sm.lastLog.serversLog
	for serverPort, rpm := range lastMetrics {
		bytes, err := json.Marshal(dto.SendServerMetricsResponse{
			Rpm: rpm,
		})
		if err != nil {
			return err
		}
		if _, ok := sm.serverPortsStatistic.Load(serverPort); ok {
			if err := sm.producer.SendMessage(fmt.Sprintf("metrics-%d", serverPort), bytes); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sm *ServerMetrics) AddServer(serverPort uint16) error {
	if _, ok := sm.serverPortsStatistic.Load(serverPort); ok {
		return nil
	}
	if err := sm.producer.CreateExchange(fmt.Sprintf("metrics-%d", serverPort)); err != nil {
		return err
	}
	sm.log.Infof("metrics have created exchange %s", fmt.Sprintf("metrics-%d", serverPort))
	sm.serverPortsStatistic.LoadOrStore(serverPort, 0)
	return nil
}

func (sm *ServerMetrics) Update(serverPort uint16) error {
	if _, ok := sm.serverPortsStatistic.Load(serverPort); !ok {
		if err := sm.AddServer(serverPort); err != nil {
			return err
		}
	}
	v, _ := sm.serverPortsStatistic.Load(serverPort)
	current := v.(int)
	sm.serverPortsStatistic.LoadOrStore(serverPort, current+1)
	return nil
}

func (sm *ServerMetrics) Delete(serverPort uint16) error {
	if err := sm.producer.DeleteExchange(fmt.Sprintf("metrics-%d", serverPort)); err != nil {
		return err
	}
	sm.log.Infof("metrics have deleted exchange %s", fmt.Sprintf("metrics-%d", serverPort))
	sm.serverPortsStatistic.Delete(serverPort)
	return nil
}
