package session

import (
	"altron/common/models"
	"altron/core/generated/spec"
	"altron/core/repositories/ent"
	"sort"
)

func PresentSessionToSpec(sessionEnt *ent.Session) *spec.Session {
	packets := make([]spec.Packet, 0)
	for _, packetEnt := range sessionEnt.Edges.Packets {
		packets = append(packets, spec.Packet{
			SentAt:    packetEnt.SentAt,
			Payload:   packetEnt.Payload,
			IsRequest: packetEnt.IsRequest,
		})
	}
	sort.Slice(packets, func(i, j int) bool {
		return packets[i].SentAt.Before(packets[j].SentAt)
	})
	return &spec.Session{
		Id:                  sessionEnt.ID.String(),
		SentAt:              sessionEnt.SentAt,
		Protocol:            sessionEnt.Protocol,
		Ttl:                 sessionEnt.TTL,
		ServerPort:          uint16(sessionEnt.ServerPort),
		ClientHost:          sessionEnt.ClientHost,
		PacketsCount:        sessionEnt.PacketsCount,
		AverageResponseTime: sessionEnt.AverageResponseTime,
		RequestsNumber:      sessionEnt.RequestsNumber,
		Packets:             packets,
		MatchedFilters:      presentSessionFiltersToSpec(sessionEnt.Edges.SessionFilters),
	}
}

func PresentSessionToModel(sessionEnt *ent.Session) *models.Session {
	packets := make([]*models.Packet, 0)
	for _, packetEnt := range sessionEnt.Edges.Packets {
		packets = append(packets, &models.Packet{
			SentAt:    packetEnt.SentAt,
			Payload:   packetEnt.Payload,
			IsRequest: packetEnt.IsRequest,
		})
	}
	sort.Slice(packets, func(i, j int) bool {
		return packets[i].SentAt.Before(packets[j].SentAt)
	})
	return &models.Session{
		ID:                  sessionEnt.ID,
		SentAt:              sessionEnt.SentAt,
		Protocol:            sessionEnt.Protocol,
		TTL:                 sessionEnt.TTL,
		ServerPort:          uint16(sessionEnt.ServerPort),
		ClientHost:          sessionEnt.ClientHost,
		PacketsCount:        sessionEnt.PacketsCount,
		AverageResponseTime: sessionEnt.AverageResponseTime,
		RequestsNumber:      sessionEnt.RequestsNumber,
		Packets:             packets,
		MatchedFilters:      presentSessionFiltersToModels(sessionEnt.Edges.SessionFilters),
	}
}

func PresentEmptySession(sessionEnt *ent.Session) *spec.Session {
	return &spec.Session{
		Id:                  sessionEnt.ID.String(),
		Ttl:                 sessionEnt.TTL,
		Protocol:            sessionEnt.Protocol,
		PacketsCount:        sessionEnt.PacketsCount,
		ServerPort:          uint16(sessionEnt.ServerPort),
		ClientHost:          sessionEnt.ClientHost,
		SentAt:              sessionEnt.SentAt,
		ClientUserAgent:     &sessionEnt.ClientUserAgent,
		AverageResponseTime: sessionEnt.AverageResponseTime,
		RequestsNumber:      sessionEnt.RequestsNumber,
		MatchedFilters:      presentSessionFiltersToSpec(sessionEnt.Edges.SessionFilters),
	}
}

func PresentEmptySessions(sessionsEnt []*ent.Session) []spec.Session {
	sessions := make([]spec.Session, 0)

	for _, sessionEnt := range sessionsEnt {
		sessions = append(sessions, *PresentEmptySession(sessionEnt))
	}
	return sessions
}

func PresentEmptySessionToModel(sessionEnt *ent.Session) *models.Session {
	return &models.Session{
		ID:                  sessionEnt.ID,
		TTL:                 sessionEnt.TTL,
		Protocol:            sessionEnt.Protocol,
		PacketsCount:        sessionEnt.PacketsCount,
		ServerPort:          uint16(sessionEnt.ServerPort),
		ClientHost:          sessionEnt.ClientHost,
		SentAt:              sessionEnt.SentAt,
		ClientUserAgent:     &sessionEnt.ClientUserAgent,
		AverageResponseTime: sessionEnt.AverageResponseTime,
		RequestsNumber:      sessionEnt.RequestsNumber,
		MatchedFilters:      presentSessionFiltersToModels(sessionEnt.Edges.SessionFilters),
	}
}

func presentSessionFiltersToSpec(sessionFiltersEnt []*ent.SessionFilter) []spec.SessionFilter {
	sessionFilters := make([]spec.SessionFilter, 0, len(sessionFiltersEnt))
	for _, sessionFilterEnt := range sessionFiltersEnt {
		filterEnt := sessionFilterEnt.Edges.Filter
		var regex *string
		var ttl *uint8
		var totalPackets *int
		if len(filterEnt.Regex) > 0 {
			regex = &filterEnt.Regex
		}
		if filterEnt.TTL != 0 {
			ttl = &filterEnt.TTL
		}
		if filterEnt.TotalPackets != 0 {
			totalPackets = &filterEnt.TotalPackets
		}
		var serviceID *string
		serviceEnt, err := filterEnt.Edges.ServiceOrErr()
		if err == nil {
			uuidStr := serviceEnt.ID.String()
			serviceID = &uuidStr
		}
		matchedPackets := make([]int, 0)
		for _, matchedPacket := range sessionFilterEnt.Edges.MatchedPackets {
			matchedPackets = append(matchedPackets, matchedPacket.PacketIdx)
		}
		sessionFilters = append(sessionFilters, spec.SessionFilter{
			Filter: spec.Filter{
				Id:           filterEnt.ID.String(),
				Name:         filterEnt.Name,
				ServiceId:    serviceID,
				Regex:        regex,
				Ttl:          ttl,
				TotalPackets: totalPackets,
				Color:        filterEnt.Color,
				InRequest:    filterEnt.InRequest,
				InResponse:   filterEnt.InResponse,
				IsBlocking:   filterEnt.IsBlocking,
			},
			MatchesCount: sessionFilterEnt.MatchesCount,
			MatchedPackets: matchedPackets,
			SessionID:    sessionFilterEnt.SessionID.String(),
		})
	}
	return sessionFilters
}

func presentSessionFiltersToModels(sessionFiltersEnt []*ent.SessionFilter) []*models.SessionFilter {
	sessionFilters := make([]*models.SessionFilter, 0, len(sessionFiltersEnt))
	for _, sessionFilterEnt := range sessionFiltersEnt {
		filterEnt := sessionFilterEnt.Edges.Filter
		var regex *string
		var ttl *uint8
		var totalPackets *int
		if len(filterEnt.Regex) > 0 {
			regex = &filterEnt.Regex
		}
		if filterEnt.TTL != 0 {
			ttl = &filterEnt.TTL
		}
		if filterEnt.TotalPackets != 0 {
			totalPackets = &filterEnt.TotalPackets
		}

		matchedPackets := make([]int, 0)
		matchedPacketsEnt := sessionFilterEnt.Edges.MatchedPackets
		for _, matchedPacketEnt := range matchedPacketsEnt {
			matchedPackets = append(matchedPackets, matchedPacketEnt.PacketIdx)
		}

		serviceEnt := filterEnt.Edges.Service
		sessionFilters = append(sessionFilters, &models.SessionFilter{
			Filter: models.Filter{
				ID:           filterEnt.ID,
				Name:         filterEnt.Name,
				ServiceID:    &serviceEnt.ID,
				Regex:        regex,
				TTL:          ttl,
				TotalPackets: totalPackets,
				Color:        filterEnt.Color,
				InRequest:    filterEnt.InRequest,
				InResponse:   filterEnt.InResponse,
				IsBlocking:   filterEnt.IsBlocking,
			},
			MatchesCount:   sessionFilterEnt.MatchesCount,
			MatchedPackets: matchedPackets,
			SessionID:      sessionFilterEnt.SessionID,
		})
	}
	return sessionFilters
}
