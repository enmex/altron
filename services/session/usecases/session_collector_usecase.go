package usecases

import (
	commonDto "altron/common/dto"
	commonHandlers "altron/common/handlers"
	common "altron/common/models"
	"altron/config"
	"altron/pkg/amqp"
	mq "altron/pkg/amqp"

	"altron/pkg/scheduler"
	"altron/session/interfaces"
	"altron/session/metrics"
	"altron/session/models"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var _ interfaces.SessionCollectorUseCase = (*SessionCollectorUseCase)(nil)

type SessionCollectorUseCase struct {
	appCfg          *config.AppConfig
	log             *logrus.Logger
	tree            *models.SessionTree
	producer        *mq.Producer
	cleanerProducer *mq.Producer
	analyzerHandler *commonHandlers.AnalyzerHandler
}

func NewSessionCollectorUseCase(
	log *logrus.Logger,
	appCfg *config.AppConfig,
	client *amqp.Client,
	analyzerHandler *commonHandlers.AnalyzerHandler,
	serverMetrics *metrics.ServerMetrics,
) (*SessionCollectorUseCase, error) {
	producer, err := amqp.NewProducer(client)
	if err != nil {
		return nil, err
	}
	cleanerProducer, err := amqp.NewProducer(client)
	if err != nil {
		return nil, err
	}

	return &SessionCollectorUseCase{
		appCfg:          appCfg,
		log:             log,
		tree:            models.NewSessionTree(appCfg.TcpStreamTimeout),
		producer:        producer,
		cleanerProducer: cleanerProducer,
		analyzerHandler: analyzerHandler,
	}, nil
}

func (u *SessionCollectorUseCase) StartSessionCleaner(ctx context.Context) *scheduler.Scheduler {
	cleaner := scheduler.NewScheduler(time.Minute, func() error {
		defer u.log.Infoln(u.TreeString())
		u.log.Infoln("session cleaner has started cleaning...")
		timeoutSessionsCount := 0
		expiredClients := make(map[*common.Session]*string, 0)
		for _, serverPort := range u.tree.Servers() {
			ifaces, err := u.tree.Interfaces(serverPort)
			if err != nil {
				return err
			}
			for _, iface := range ifaces {
				clientHosts, err := u.tree.ClientHosts(serverPort, iface)
				if err != nil {
					continue
				}
				for _, clientHost := range clientHosts {
					session, fileName, err := u.tree.ClientSession(serverPort, iface, clientHost)
					if err != nil {
						continue
					}
					oldPacket := session.Packets[0]
					var timeout time.Duration
					if strings.EqualFold(session.Protocol, "tcp") {
						timeout = u.appCfg.TcpStreamTimeout
					} else {
						timeout = u.appCfg.UdpStreamTimeout
					}

					if time.Since(oldPacket.SentAt) > timeout {
						expiredClients[session] = fileName
						timeoutSessionsCount++
					}
				}
			}
		}

		for session, fileName := range expiredClients {
			processedSession, err := u.processSession(ctx, session, fileName)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(session)
			if err != nil {
				return err
			}
			var exchangeName string
			if fileName != nil {
				exchangeName = *fileName
			} else {
				exchangeName = fmt.Sprint(processedSession.ServerPort)
			}

			if err := u.producer.SendMessage(exchangeName, bytes); err != nil {
				return err
			}
			if err := u.tree.DeleteClient(processedSession.Iface, processedSession.ServerPort, processedSession.ClientHost); err != nil {
				return err
			}
		}

		if timeoutSessionsCount > 0 {
			u.log.Infof("Session cleaner has cleared %d expired session(s)", timeoutSessionsCount)
		}
		return nil
	})
	cleaner.Start()

	return cleaner
}

func (u *SessionCollectorUseCase) AddPacket(ctx context.Context, payload commonDto.SendPacketToQueueRequest) error {
	if !u.tree.ServerExists(payload.ServerPort) && payload.FileName == nil {
		return nil
	} else if payload.FileName != nil {
		u.tree.AddServer(payload.ServerPort)
	}
	if !u.tree.PortInterfaceExists(payload.ServerPort, payload.Interface) {
		if err := u.tree.AddInterface(payload.ServerPort, payload.Interface); err != nil {
			return err
		}
	}
	exists, err := u.tree.ClientExists(payload.ServerPort, payload.Interface, payload.ClientHost)
	if err != nil {
		return err
	}
	if !exists {
		userAgent, err := u.getClientUserAgent(payload.Packet.Payload)
		if err != nil {
			return err
		}
		if err := u.tree.AddClientHost(
			payload.ServerPort,
			payload.Interface,
			payload.ClientHost,
			payload.Protocol,
			payload.Packet.SentAt,
			payload.TTL,
			payload.FileName,
			userAgent,
		); err != nil {
			return err
		}
	}
	if payload.IsLast {
		return u.sendSession(ctx, payload.ServerPort, payload.Interface, payload.ClientHost)
	}

	return u.tree.AddPacket(payload.ServerPort, payload.Interface, payload.ClientHost, &payload.Packet)
}

func (u *SessionCollectorUseCase) DeleteServer(serverPort uint16) error {
	u.tree.DeleteServer(serverPort)
	return u.producer.DeleteQueue(fmt.Sprint(serverPort))
}

func (u *SessionCollectorUseCase) GetServers() []uint16 {
	return u.tree.Servers()
}

func (u *SessionCollectorUseCase) ServerExists(serverPort uint16) bool {
	for _, port := range u.tree.Servers() {
		if serverPort == port {
			return true
		}
	}
	return false
}

func (u *SessionCollectorUseCase) AddServer(serverPort uint16) error {
	u.tree.AddServer(serverPort)
	return u.producer.CreateExchange(fmt.Sprint(serverPort))
}

func (u *SessionCollectorUseCase) sendSession(ctx context.Context, serverPort uint16, iface string, clientHost string) error {
	session, fileName, err := u.tree.ClientSession(serverPort, iface, clientHost)
	if err != nil {
		return fmt.Errorf("%s: %d %s", err.Error(), serverPort, clientHost)
	}
	if len(session.Packets) > 1 {
		session, err := u.processSession(ctx, session, fileName)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(session)
		if err != nil {
			return err
		}
		var exchangeName string
		if fileName != nil {
			exchangeName = *fileName
		} else {
			exchangeName = fmt.Sprint(serverPort)
		}

		if err := u.producer.SendMessage(exchangeName, bytes); err != nil {
			return err
		}
	}
	if err := u.tree.DeleteClient(iface, serverPort, clientHost); err != nil {
		return err
	}
	return nil
}

func (u *SessionCollectorUseCase) processSession(ctx context.Context, session *common.Session, fileName *string) (*common.Session, error) {
	mergedPackets, err := u.mergeAdjacentPackets(session.Packets)
	if err != nil {
		return nil, err
	}
	session.Packets = mergedPackets
	session.RequestsNumber = 0
	var timestampsDiffSum int64 = 0
	for idx, p := range mergedPackets {
		if p.IsRequest {
			session.RequestsNumber++
		}
		if idx%2 == 1 {
			timestampsDiffSum += session.Packets[idx].SentAt.Sub(session.Packets[idx-1].SentAt).Milliseconds()
		}
	}

	session.AverageResponseTime = float64(timestampsDiffSum) / float64(2*len(session.Packets))
	session.PacketsCount = len(mergedPackets)
	if fileName == nil {
		analyzerMatches, err := u.analyzerHandler.ServeSession(ctx, session)
		if err != nil {
			return nil, err
		}
		session.AnalyzerMatches = analyzerMatches
	}
	return session, nil
}

func (u *SessionCollectorUseCase) mergeAdjacentPackets(packets []*common.Packet) ([]*common.Packet, error) {
	result := make([]*common.Packet, 0)
	result = append(result, packets[0])
	k := 0

	if len(packets) >= 1000 {
		return packets, nil
	}

	for i := 1; i < len(packets); i++ {
		if strings.EqualFold(result[k].Payload, packets[i].Payload) {
			continue
		}
		if packets[i].IsRequest == result[k].IsRequest {
			decodedPacket, err := base64.StdEncoding.DecodeString(result[k].Payload)
			if err != nil {
				return nil, err
			}
			decodedPart, err := base64.StdEncoding.DecodeString(packets[i].Payload)
			if err != nil {
				return nil, err
			}
			result[k].Payload = base64.StdEncoding.EncodeToString(append(decodedPacket, decodedPart...))
		} else {
			result = append(result, packets[i])
			k++
		}
	}
	return result, nil
}

func (u *SessionCollectorUseCase) getClientUserAgent(payload string) (*string, error) {
	decodedPacket, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(decodedPacket))))
	if err != nil {
		return nil, nil
	}
	userAgent := req.Header.Get("User-Agent")
	return &userAgent, nil
}

func (u *SessionCollectorUseCase) TreeString() string {
	return u.tree.String()
}
