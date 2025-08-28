package usecases

import (
	"altron/pkg/plugin"
	"altron/plugin/generated/spec"
	"altron/plugin/interfaces"
	"encoding/base64"
	"strings"

	"github.com/sirupsen/logrus"
)

var _ interfaces.PluginUseCase = (*PluginUseCase)(nil)

type PluginUseCase struct {
	log           *logrus.Logger
	pluginManager *plugin.PluginManager[interfaces.PluginInterface]
}

func NewPluginUseCase(log *logrus.Logger, pluginManager *plugin.PluginManager[interfaces.PluginInterface]) *PluginUseCase {
	return &PluginUseCase{
		log:           log,
		pluginManager: pluginManager,
	}
}

func (u *PluginUseCase) ProcessSession(request *spec.ProcessSessionRequest) (*spec.ProcessSessionResponse, error) {
	plugins := u.pluginManager.GetPluginsByNames(request.Plugins)
	packets := request.Session.Packets

	rawSession := make([][]byte, 0, len(packets))
	for _, packet := range packets {
		decodedPayload, err := base64.StdEncoding.DecodeString(packet.Payload)
		if err != nil {
			return nil, err
		}
		rawSession = append(rawSession, decodedPayload)
	}

	var err error
	for _, plugin := range plugins {
		rawSession, err = plugin.Process(rawSession)
		if err != nil {
			if strings.Contains(err.Error(), "gzip: invalid header") {
				u.log.Infoln(request.Session.Packets)
			}
			return nil, err
		}
	}
	parsedPackets := make([]spec.Packet, 0, len(rawSession))
	for idx, rawPacketPayload := range rawSession {
		encodedPacketPayload := base64.StdEncoding.EncodeToString(rawPacketPayload)
		parsedPackets = append(parsedPackets, spec.Packet{
			Payload:   encodedPacketPayload,
			IsRequest: packets[idx].IsRequest,
			SentAt:    packets[idx].SentAt,
		})
	}

	session := request.Session
	session.Packets = parsedPackets

	return &spec.ProcessSessionResponse{
		Session: session,
	}, nil
}

func (u *PluginUseCase) GetAllPlugins() *spec.GetAllPluginsResponse {
	return &spec.GetAllPluginsResponse{
		Plugins: u.pluginManager.GetPluginNames(),
	}
}
