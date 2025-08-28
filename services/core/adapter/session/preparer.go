package session

import (
	"altron/core/generated/spec"
	"bufio"
	"encoding/base64"
	"net/http"
	"strings"

	common "altron/common/models"

	"github.com/google/uuid"
)

func PrepareCreateSessionRequest(request *spec.CreateSessionRequest) (*common.Session, error) {
	sessionSpec := request.Session
	packets := make([]*common.Packet, 0)
	for _, packetSpec := range sessionSpec.Packets {
		packets = append(packets, &common.Packet{
			SentAt:    packetSpec.SentAt,
			Payload:   packetSpec.Payload,
			IsRequest: packetSpec.IsRequest,
		})
	}
	filters := make([]*common.SessionFilter, 0, len(sessionSpec.MatchedFilters))
	for _, filterSpec := range sessionSpec.MatchedFilters {
		filters = append(filters, &common.SessionFilter{
			Filter: common.Filter{
				ID:         uuid.MustParse(filterSpec.Id),
				Name:       filterSpec.Name,
				Color:      filterSpec.Color,
				InRequest:  filterSpec.InRequest,
				InResponse: filterSpec.InResponse,
				IsBlocking: filterSpec.IsBlocking,
				Regex:      filterSpec.Regex,
				TTL:        filterSpec.Ttl,
			},
			MatchesCount:   filterSpec.MatchesCount,
			MatchedPackets: filterSpec.MatchedPackets,
		})
	}
	var userAgent string
	decodedPacket, err := base64.StdEncoding.DecodeString(sessionSpec.Packets[0].Payload)
	if err != nil {
		return nil, err
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(decodedPacket))))
	if err == nil {
		userAgent = req.Header.Get("User-Agent")
	}

	return &common.Session{
		ID:                  uuid.MustParse(sessionSpec.Id),
		SentAt:              sessionSpec.SentAt,
		Protocol:            sessionSpec.Protocol,
		ServerPort:          sessionSpec.ServerPort,
		ClientHost:          sessionSpec.ClientHost,
		TTL:                 sessionSpec.Ttl,
		PacketsCount:        sessionSpec.PacketsCount,
		AverageResponseTime: sessionSpec.AverageResponseTime,
		RequestsNumber:      sessionSpec.RequestsNumber,
		ClientUserAgent:     &userAgent,
		Packets:             packets,
		MatchedFilters:      filters,
	}, nil
}

func PrepareAddSessionToPcapWorkspaceRequest(request *spec.AddPcapSessionRequest) (*common.Session, error) {
	sessionSpec := request.Session

	packets := make([]*common.Packet, 0)
	for _, packetSpec := range sessionSpec.Packets {
		packets = append(packets, &common.Packet{
			SentAt:    packetSpec.SentAt,
			Payload:   packetSpec.Payload,
			IsRequest: packetSpec.IsRequest,
		})
	}
	filters := make([]*common.SessionFilter, 0, len(sessionSpec.MatchedFilters))
	for _, filterSpec := range sessionSpec.MatchedFilters {
		filters = append(filters, &common.SessionFilter{
			Filter: common.Filter{
				ID:         uuid.MustParse(filterSpec.Id),
				Name:       filterSpec.Name,
				Color:      filterSpec.Color,
				InRequest:  filterSpec.InRequest,
				InResponse: filterSpec.InResponse,
				IsBlocking: filterSpec.IsBlocking,
				Regex:      filterSpec.Regex,
				TTL:        filterSpec.Ttl,
			},
			MatchesCount:   filterSpec.MatchesCount,
			MatchedPackets: filterSpec.MatchedPackets,
		})
	}
	var userAgent string
	decodedPacket, err := base64.StdEncoding.DecodeString(sessionSpec.Packets[0].Payload)
	if err != nil {
		return nil, err
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(decodedPacket))))
	if err == nil {
		userAgent = req.Header.Get("User-Agent")
	}

	return &common.Session{
		ID:                  uuid.MustParse(sessionSpec.Id),
		SentAt:              sessionSpec.SentAt,
		Protocol:            sessionSpec.Protocol,
		ServerPort:          sessionSpec.ServerPort,
		ClientHost:          sessionSpec.ClientHost,
		TTL:                 sessionSpec.Ttl,
		PacketsCount:        sessionSpec.PacketsCount,
		ClientUserAgent:     &userAgent,
		AverageResponseTime: sessionSpec.AverageResponseTime,
		RequestsNumber:      sessionSpec.RequestsNumber,
		Packets:             packets,
		MatchedFilters:      filters,
	}, nil
}

func PrepareAddSessionsToWorkspaceRequest(request *spec.AddSessionsToWorkspaceRequest) ([]*common.Session, error) {
	sessions := make([]*common.Session, 0)
	for _, sessionSpec := range request.Sessions {
		packets := make([]*common.Packet, 0)
		for _, packetSpec := range sessionSpec.Packets {
			packets = append(packets, &common.Packet{
				SentAt:    packetSpec.SentAt,
				Payload:   packetSpec.Payload,
				IsRequest: packetSpec.IsRequest,
			})
		}
		filters := make([]*common.SessionFilter, 0, len(sessionSpec.MatchedFilters))
		for _, filterSpec := range sessionSpec.MatchedFilters {
			filters = append(filters, &common.SessionFilter{
				Filter: common.Filter{
					ID:         uuid.MustParse(filterSpec.Id),
					Name:       filterSpec.Name,
					Color:      filterSpec.Color,
					InRequest:  filterSpec.InRequest,
					InResponse: filterSpec.InResponse,
					IsBlocking: filterSpec.IsBlocking,
					Regex:      filterSpec.Regex,
					TTL:        filterSpec.Ttl,
				},
				MatchesCount:   filterSpec.MatchesCount,
				MatchedPackets: filterSpec.MatchedPackets,
			})
		}
		var userAgent string
		decodedPacket, err := base64.StdEncoding.DecodeString(sessionSpec.Packets[0].Payload)
		if err != nil {
			return nil, err
		}
		req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(decodedPacket))))
		if err == nil {
			userAgent = req.Header.Get("User-Agent")
		}

		sessions = append(sessions, &common.Session{
			ID:                  uuid.MustParse(sessionSpec.Id),
			SentAt:              sessionSpec.SentAt,
			Protocol:            sessionSpec.Protocol,
			ServerPort:          sessionSpec.ServerPort,
			ClientHost:          sessionSpec.ClientHost,
			TTL:                 sessionSpec.Ttl,
			PacketsCount:        sessionSpec.PacketsCount,
			ClientUserAgent:     &userAgent,
			AverageResponseTime: sessionSpec.AverageResponseTime,
			RequestsNumber:      sessionSpec.RequestsNumber,
			Packets:             packets,
			MatchedFilters:      filters,
		})
	}
	return sessions, nil
}

func PrepareConvertSessionToExploit(request *spec.ConvertSessionToExploitRequest) *common.Session {
	packets := make([]*common.Packet, 0)
	for _, packetSpec := range request.Session.Packets {
		packets = append(packets, &common.Packet{
			SentAt:    packetSpec.SentAt,
			Payload:   packetSpec.Payload,
			IsRequest: packetSpec.IsRequest,
		})
	}
	return &common.Session{
		TTL:          request.Session.Ttl,
		Protocol:     request.Session.Protocol,
		SentAt:       request.Session.SentAt,
		ServerPort:   request.Session.ServerPort,
		ClientHost:   request.Session.ClientHost,
		PacketsCount: request.Session.PacketsCount,
		Packets:      packets,
	}
}

func PrepareToModel(sessionSpec spec.Session) *common.Session {
	session := &common.Session{
		ID:                  uuid.MustParse(sessionSpec.Id),
		TTL:                 sessionSpec.Ttl,
		Protocol:            sessionSpec.Protocol,
		PacketsCount:        sessionSpec.PacketsCount,
		ServerPort:          sessionSpec.ServerPort,
		ClientHost:          sessionSpec.ClientHost,
		SentAt:              sessionSpec.SentAt,
		ClientUserAgent:     sessionSpec.ClientUserAgent,
		AverageResponseTime: sessionSpec.AverageResponseTime,
		RequestsNumber:      sessionSpec.RequestsNumber,
		MatchedFilters:      prepareSessionFilters(sessionSpec.MatchedFilters),
	}
	return session
}

func prepareSessionFilters(sessionFiltersSpec []spec.SessionFilter) []*common.SessionFilter {
	sessionFilters := make([]*common.SessionFilter, 0, len(sessionFiltersSpec))
	for _, sessionFilterSpec := range sessionFiltersSpec {
		sessionFilters = append(sessionFilters, &common.SessionFilter{
			Filter: common.Filter{
				ID:           uuid.MustParse(sessionFilterSpec.Id),
				Name:         sessionFilterSpec.Name,
				Regex:        sessionFilterSpec.Regex,
				TTL:          sessionFilterSpec.Ttl,
				TotalPackets: sessionFilterSpec.TotalPackets,
				Color:        sessionFilterSpec.Color,
				InRequest:    sessionFilterSpec.InRequest,
				InResponse:   sessionFilterSpec.InResponse,
				IsBlocking:   sessionFilterSpec.IsBlocking,
			},
			MatchesCount:   sessionFilterSpec.MatchesCount,
			MatchedPackets: sessionFilterSpec.MatchedPackets,
			SessionID:      uuid.MustParse(sessionFilterSpec.SessionID),
		})
	}
	return sessionFilters
}
