package handlers

import (
	common "altron/common/models"
	"altron/config"
	"altron/connection/dto"
	"altron/pkg/request"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	log *logrus.Logger
	cfg *config.AppConfig
}

func NewSessionHandler(log *logrus.Logger, cfg *config.AppConfig) *SessionHandler {
	return &SessionHandler{
		log: log,
		cfg: cfg,
	}
}

func (h *SessionHandler) ServeSession(session *common.Session, serviceFilters []*common.Filter, servicePlugins []string) (*common.Session, error) {
	if len(servicePlugins) != 0 {
		processSessionResponse, err := request.Post[dto.ProcessSessionResponse](
			fmt.Sprintf("http://altron.plugin.loc:%d/plugins/process", h.cfg.AltronPluginPort),
			dto.ProcessSessionRequest{
				Session: session,
				Plugins: servicePlugins,
			},
		)
		if err != nil {
			h.log.Errorln(err)
		} else if len(processSessionResponse.Packets) > 0 {
			session.Packets = processSessionResponse.Packets
			session.PacketsCount = len(session.Packets)
		}
	}
	spam := true
	var err error
	session.MatchedFilters, spam, err = h.getSessionMatchedFilters(session, serviceFilters)
	if err != nil {
		return nil, err
	}
	if spam {
		return nil, nil
	}
	return session, nil
}

func (h *SessionHandler) getSessionMatchedFilters(session *common.Session, filters []*common.Filter) ([]*common.SessionFilter, bool, error) {
	matchedFilters := make([]*common.SessionFilter, 0, len(filters))

	for _, filter := range filters {
		match := true
		matchesCount := 0
		matchedPackets := make([]int, 0)

		if filter.TTL != nil {
			match = session.TTL == *filter.TTL
		}
		if filter.TotalPackets != nil {
			match = match && session.PacketsCount/2 == *filter.TotalPackets
		}

		if filter.Regex != nil {
			payloadMatches := false
			for idx, packet := range session.Packets {
				payload, err := base64.StdEncoding.DecodeString(packet.Payload)
				if err != nil {
					return nil, false, err
				}
				if packet.IsRequest && filter.InRequest || !packet.IsRequest && filter.InResponse {
					re := regexp.MustCompile(*filter.Regex)
					matches := re.FindAllStringIndex(string(payload), -1)
					if !payloadMatches {
						payloadMatches = len(matches) != 0
					}
					matchesCount += len(matches)
					if len(matches) > 0 {
						matchedPackets = append(matchedPackets, idx)
					}
				}
			}
			match = match && payloadMatches
		}

		if match {
			if filter.IsBlocking {
				return nil, true, nil
			}
			matchedFilters = append(matchedFilters, &common.SessionFilter{
				Filter:       *filter,
				SessionID:    session.ID,
				MatchesCount: matchesCount,
				MatchedPackets: matchedPackets,
			})
		}
	}
	return matchedFilters, false, nil
}
