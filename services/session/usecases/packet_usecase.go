package usecases

import (
	dto "altron/common/dto"
	"altron/common/models"
	"altron/config"
	"altron/pkg/packets"
	"altron/pkg/scheduler"
	sftp "altron/pkg/sftp"
	"altron/session/generated/spec"
	"altron/session/interfaces"
	"altron/session/metrics"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

var _ interfaces.PacketUseCase = (*PacketUseCase)(nil)

type ClientData struct {
	CreatedAt      time.Time
	LastTcpPackets []*layers.TCP
}

type PacketUseCase struct {
	log              *logrus.Logger
	appCfg           *config.AppConfig
	sftpClient       *sftp.Client
	sessionCollector *SessionCollectorUseCase
	logsCollector    *LogsCollectorUseCase
	pcapAnalyzer     *packets.PcapAnalyzer
	serverMetrics    *metrics.ServerMetrics
	lastTcpPackets   sync.Map //*utils.ConcurrentMap[string, *ClientData] //client port to last packets
	ports            sync.Map //*utils.ConcurrentMap[uint16, bool]
}

func NewPacketUseCase(
	log *logrus.Logger,
	appCfg *config.AppConfig,
	sftpClient *sftp.Client,
	sessionCollector *SessionCollectorUseCase,
	logsCollector *LogsCollectorUseCase,
	pcapAnalyzer *packets.PcapAnalyzer,
	serverMetrics *metrics.ServerMetrics,
) (*PacketUseCase, error) {
	return &PacketUseCase{
		log:              log,
		appCfg:           appCfg,
		sftpClient:       sftpClient,
		sessionCollector: sessionCollector,
		logsCollector:    logsCollector,
		pcapAnalyzer:     pcapAnalyzer,
		serverMetrics:    serverMetrics,
		lastTcpPackets:   sync.Map{}, //utils.NewConcurrentMap[string, *ClientData](),
		ports:            sync.Map{}, //utils.NewConcurrentMap[uint16, bool](),
	}, nil
}

func (u *PacketUseCase) AddPcapPorts(ctx context.Context, request *spec.CreatePortsRequest) error {
	for _, req := range request.Ports {
		u.ports.LoadOrStore(req.Port, true)
		if err := u.sessionCollector.AddServer(req.Port); err != nil {
			return err
		}
		if err := u.serverMetrics.AddServer(req.Port); err != nil {
			return err
		}
		if req.ContainerID != nil {
			if err := u.logsCollector.AddLogsContainer(*req.ContainerID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *PacketUseCase) DeletePcapPort(ctx context.Context, request *spec.DeletePortRequest) error {
	u.ports.Delete(request.Port)
	if err := u.sessionCollector.DeleteServer(request.Port); err != nil {
		return err
	}
	if err := u.serverMetrics.Delete(request.Port); err != nil {
		return err
	}
	if request.ContainerID != nil {
		if err := u.logsCollector.AddLogsContainer(*request.ContainerID); err != nil {
			return err
		}
	}
	return nil
}

func (u *PacketUseCase) StartMetricsMeasure() *scheduler.Scheduler {
	measureScheduler := scheduler.NewScheduler(time.Minute, func() error {
		u.serverMetrics.Checkpoint()
		return nil
	})
	measureScheduler.Start()

	return measureScheduler
}

func (u *PacketUseCase) StartSendMetrics() *scheduler.Scheduler {
	senderScheduler := scheduler.NewScheduler(time.Second, func() error {
		return u.serverMetrics.SendMetrics()
	})
	senderScheduler.Start()

	return senderScheduler
}

func (u *PacketUseCase) ProducePacket(ctx context.Context, packet gopacket.Packet, fileName *string) {
	if packet.TransportLayer() == nil {
		return
	}
	transportFlow := packet.TransportLayer().TransportFlow()
	srcPort, _ := strconv.Atoi(transportFlow.Src().String())
	dstPort, _ := strconv.Atoi(transportFlow.Dst().String())

	if fileName != nil {
		_, srcOk := u.ports.Load(uint16(srcPort))
		_, dstOk := u.ports.Load(uint16(dstPort))
		if !srcOk && !dstOk {
			u.ports.LoadOrStore(uint16(dstPort), true)
		}
	}

	var serverHost net.IP
	var clientHost string
	var serverPort uint16

	isRequest := false
	var ttl uint8
	ipv4Layer, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	if ok {
		ttl = ipv4Layer.TTL
		if _, ok := u.ports.Load(uint16(srcPort)); ok {
			clientHost = fmt.Sprintf("%s:%d", ipv4Layer.DstIP.String(), dstPort)
			serverPort = uint16(srcPort)
			serverHost = ipv4Layer.SrcIP
		} else {
			clientHost = fmt.Sprintf("%s:%d", ipv4Layer.SrcIP.String(), srcPort)
			serverPort = uint16(dstPort)
			serverHost = ipv4Layer.DstIP
			isRequest = true
		}
	} else {
		ipv6Layer, ok := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		if !ok {
			return
		}
		ttl = ipv6Layer.HopLimit
		if _, ok := u.ports.Load(uint16(srcPort)); ok {
			clientHost = fmt.Sprintf("%s:%d", ipv6Layer.DstIP.String(), dstPort)
			serverPort = uint16(srcPort)
			serverHost = ipv6Layer.SrcIP
		} else {
			clientHost = fmt.Sprintf("%s:%d", ipv6Layer.SrcIP.String(), srcPort)
			serverPort = uint16(dstPort)
			serverHost = ipv6Layer.DstIP
			isRequest = true
		}
	}

	var payload []byte

	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		payload = applicationLayer.Payload()
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payload)

	var protocol string
	isLast := false

	tcpLayer, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
	if ok {
		isLast = u.isLastTcpPacket(clientHost, tcpLayer)
		protocol = "tcp"
	} else {
		protocol = "udp"
	}

	if len(payload) != 0 || isLast {
		if err := u.sessionCollector.AddPacket(ctx, dto.SendPacketToQueueRequest{
			FileName:   fileName,
			Interface:  serverHost.String(),
			ClientHost: clientHost,
			ServerPort: serverPort,
			TTL:        ttl,
			Protocol:   protocol,
			Packet: models.Packet{
				SentAt:    packet.Metadata().Timestamp,
				Payload:   encodedPayload,
				IsRequest: isRequest,
			},
			IsLast: isLast,
		}); err != nil {
			u.log.Errorln(err)
		}
		if fileName == nil {
			if err := u.serverMetrics.Update(serverPort); err != nil {
				u.log.Errorln(err)
			}
		}
	}
}

func (u *PacketUseCase) isLastTcpPacket(clientHost string, tcpLayer *layers.TCP) bool {
	if _, ok := u.lastTcpPackets.Load(clientHost); !ok {
		u.lastTcpPackets.LoadOrStore(clientHost, &ClientData{
			CreatedAt:      time.Now(),
			LastTcpPackets: make([]*layers.TCP, 0),
		})
	}
	if tcpLayer.RST {
		u.lastTcpPackets.Delete(clientHost)
		return true
	}
	if tcpLayer.ACK && tcpLayer.FIN {
		v, _ := u.lastTcpPackets.Load(clientHost)
		clientData := v.(*ClientData)
		lastPackets := clientData.LastTcpPackets
		lastPackets = append(lastPackets, tcpLayer)

		if len(lastPackets) == 2 {
			u.lastTcpPackets.Delete(clientHost)
			return true
		}
		u.lastTcpPackets.LoadOrStore(clientHost, &ClientData{
			CreatedAt:      clientData.CreatedAt,
			LastTcpPackets: lastPackets,
		})
	}
	return false
}

func (u *PacketUseCase) ImportFile(ctx context.Context, request *spec.UploadPcapRequest) error {
	filepath := fmt.Sprintf("/files/%s", request.FileName)

	if err := u.sftpClient.Download(filepath); err != nil {
		return err
	}
	defer os.Remove(filepath)
	defer u.sftpClient.Delete(filepath)

	handler, err := u.pcapAnalyzer.PcapFileHandler(filepath)
	if err != nil {
		return err
	}
	defer handler.Close()

	packetsChan := u.pcapAnalyzer.PacketSource(handler)
	for {
		packet, ok := <-packetsChan
		if !ok {
			u.log.Infof("pcap file %s has been read and processed", filepath)
			return nil
		}
		u.ProducePacket(ctx, packet, &request.FileName)
	}
}

func (u *PacketUseCase) MonitorDumps(ctx context.Context) error {
	for {
		time.Sleep(time.Second)
		list, err := u.sftpClient.List("pcaps")
		if err != nil {
			return err
		}
		if len(list) == 0 {
			continue
		}
		wg := &sync.WaitGroup{}
		wg.Add(len(list))
		for _, entry := range list {
			path := "/pcaps/" + entry
			u.log.Infoln("found pcap file", path)
			if err := u.sftpClient.Download(path); err != nil {
				wg.Done()
				u.log.Errorln(fmt.Sprintf("error downloading entry %s: %v", entry, err))
				continue
			}

			go u.processPcapDump(ctx, path, wg)
		}
		wg.Wait()
	}
}

func (u *PacketUseCase) ClearExpiredClients(ctx context.Context) *scheduler.Scheduler {
	clearExpiredClientsScheduler := scheduler.NewScheduler(time.Minute, func() error {
		expiredNumber := 0
		clients := make([]string, 0)
		u.lastTcpPackets.Range(func(key, value any) bool {
			clients = append(clients, key.(string))
			return true
		})
		for _, client := range clients {
			v, _ := u.lastTcpPackets.Load(client)
			clientData := v.(*ClientData)
			if time.Now().Sub(clientData.CreatedAt) > u.appCfg.TcpStreamTimeout {
				expiredNumber++
				u.lastTcpPackets.Delete(client)
			}
		}
		u.log.Infof("cleared %d clients, current number of clients to handle: %d", expiredNumber, len(clients))
		u.log.Infoln(u.sessionCollector.TreeString())
		return nil
	})
	clearExpiredClientsScheduler.Start()

	return clearExpiredClientsScheduler
}

func (u *PacketUseCase) processPcapDump(ctx context.Context, pcapPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer os.Remove(pcapPath)
	defer u.sftpClient.Delete(pcapPath)
	// pcapPath, err := u.ungzipPcap(gzipPath)
	// if err != nil {
	// 	u.log.Errorln("pcap ungzip error:", err)
	// 	return
	// }
	// defer os.Remove(pcapPath)

	handler, err := u.pcapAnalyzer.PcapFileHandler(pcapPath)
	if err != nil {
		u.log.Errorln(err)
		return
	}
	defer handler.Close()

	packetsChan := u.pcapAnalyzer.PacketSource(handler)

	for {
		select {
		case packet, ok := <-packetsChan:
			if !ok {
				u.log.Infof("pcap file %s has been read and processed", pcapPath)
				u.log.Infoln(u.sessionCollector.TreeString())
				return
			}
			u.ProducePacket(ctx, packet, nil)
		case <-time.After(3 * time.Second):
			u.log.Infof("pcap file %s timed out, no packets in channel", pcapPath)
			return
		}
	}
}

func (u *PacketUseCase) ungzipPcap(gzipPath string) (string, error) {
	pcapGz, err := os.Open(gzipPath)
	if err != nil {
		return "", err
	}
	defer pcapGz.Close()

	gzipReader, err := gzip.NewReader(pcapGz)
	if err != nil {
		u.log.Errorln("Unable to create gzip reader")
		return "", err
	}
	defer gzipReader.Close()

	pcapPath := fmt.Sprintf("/pcaps/%s.pcap", uuid.New())

	outputPcap, err := os.Create(pcapPath)
	if err != nil {
		u.log.Errorln("Could not create output pcap file", pcapPath)
		return "", err
	}
	defer outputPcap.Close()

	if _, err := io.Copy(outputPcap, gzipReader); err != nil {
		u.log.Errorln("Could not copy output pcap file")
		return "", err
	}

	return pcapPath, nil
}
