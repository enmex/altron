package usecases

import (
	"altron/common/dto"
	commonHandlers "altron/common/handlers"
	"altron/common/models"
	"altron/core/adapter/packet"
	"altron/core/adapter/session"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/repositories/ent"
	"context"
	"encoding/base64"
	"regexp"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var _ interfaces.SessionUseCase = (*SessionUseCase)(nil)

type SessionUseCase struct {
	log                *logrus.Logger
	sessionRepo        interfaces.SessionRepository
	cartRepo           interfaces.CartRepository
	analyzerHandler    *commonHandlers.AnalyzerHandler
}

func NewSessionUseCase(
	log *logrus.Logger,
	sessionRepo interfaces.SessionRepository,
	cartRepo interfaces.CartRepository,
	analyzerHandler *commonHandlers.AnalyzerHandler,
) *SessionUseCase {
	return &SessionUseCase{
		log:                log,
		sessionRepo:        sessionRepo,
		cartRepo:           cartRepo,
		analyzerHandler:    analyzerHandler,
	}
}

func (u *SessionUseCase) AddSessions(ctx context.Context, workspaceID uuid.UUID, request *spec.AddSessionsToWorkspaceRequest) error {
	preparedSession, err := session.PrepareAddSessionsToWorkspaceRequest(request)
	if err != nil {
		return err
	}
	return u.sessionRepo.CreateSessions(ctx, workspaceID, preparedSession)
}

func (u *SessionUseCase) GetPaginatedSessions(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) (*spec.GetSessionsResponse, error) {
	sessionsEnt, err := u.sessionRepo.GetSessionsByWorkspace(ctx, workspaceID, nil, paginationIndex)
	if err != nil {
		return nil, err
	}
	return &spec.GetSessionsResponse{
		Sessions: session.PresentEmptySessions(sessionsEnt),
	}, nil
}

func (u *SessionUseCase) ClearSessions(ctx context.Context, workspaceID uuid.UUID) error {
	return u.sessionRepo.DeleteAllSessions(ctx, workspaceID)
}

func (u *SessionUseCase) SearchSessions(ctx context.Context, workspaceID uuid.UUID, paginationIndex int, request *spec.SearchSessionsRequest) (*spec.SearchSessionsResponse, error) {
	resultSessions := make([]spec.Session, 0)

	var filterUuid *uuid.UUID
	if request.FilterID != nil {
		uuid := uuid.MustParse(*request.FilterID)
		filterUuid = &uuid
	}
	sessionCount, err := u.sessionRepo.CountWorkspaceSessions(ctx, workspaceID, filterUuid)
	if err != nil {
		return nil, err
	}

	paginationCount := 1 + sessionCount/100
	wg := sync.WaitGroup{}
	wg.Add(paginationCount)

	for paginationIndex := 0; paginationIndex < paginationCount; paginationIndex++ {
		sessions, err := u.sessionRepo.GetSessionsByWorkspace(ctx, workspaceID, filterUuid, paginationIndex)
		if err != nil {
			return nil, err
		}
		go func() {
			defer wg.Done()
			for _, s := range sessions {
				if request.SearchValue != nil {
					packets, err := u.sessionRepo.GetSessionPackets(ctx, s.ID)
					if err != nil {
						u.log.Errorln(err)
						continue
					}
					for _, packet := range packets {
						payload, err := base64.StdEncoding.DecodeString(packet.Payload)
						if err != nil {
							u.log.Errorln(err)
							continue
						}
						matches, err := regexp.Match(*request.SearchValue, payload)
						if err != nil {
							u.log.Errorln(err)
							continue
						}
						if matches {
							resultSessions = append(resultSessions, *session.PresentEmptySession(s))
							break
						}
					}
				} else {
					resultSessions = append(resultSessions, *session.PresentEmptySession(s))
				}
			}
		}()
	}
	wg.Wait()
	resultSlice := make([]spec.Session, 0, len(resultSessions))

	if len(request.SelectedCharacteristics.AdditionalProperties) > 0 {
		characteristicsPayload := make(map[string][]dto.SelectedCharacteristicDto, 0)
		for characteristicName, characteristics := range request.SelectedCharacteristics.AdditionalProperties {
			characteristicsDto := make([]dto.SelectedCharacteristicDto, 0, len(characteristics))
			for _, ch := range characteristics {
				characteristicsDto = append(characteristicsDto, dto.SelectedCharacteristicDto{
					Value:    ch.Value,
					Selected: ch.Selected,
					Blocked:  ch.Blocked,
				})
			}
			characteristicsPayload[characteristicName] = characteristicsDto
		}
		for _, s := range resultSessions {
			if u.analyzerHandler.MatchesSelectedCharacteristic(
				ctx,
				session.PrepareToModel(s),
				characteristicsPayload,
			) {
				resultSlice = append(resultSlice, s)
			}
		}
	} else {
		resultSlice = resultSessions
	}
	if len(resultSlice) > 100*(paginationIndex+1) {
		resultSlice = resultSlice[100*paginationIndex : 100*(paginationIndex+1)]
	}

	sort.Slice(resultSlice, func(i, j int) bool {
		return resultSlice[i].SentAt.Before(resultSlice[j].SentAt)
	})

	return &spec.SearchSessionsResponse{
		Sessions: resultSlice,
	}, nil
}

func (u *SessionUseCase) GetSession(ctx context.Context, sessionID uuid.UUID) (*spec.GetSessionResponse, error) {
	sessionEnt, err := u.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return &spec.GetSessionResponse{
		Session: *session.PresentSessionToSpec(sessionEnt),
	}, nil
}

func (u *SessionUseCase) CreateSession(ctx context.Context, userID uuid.UUID, request *spec.CreateSessionRequest) error {
	session, err := session.PrepareCreateSessionRequest(request)
	if err != nil {
		return err
	}
	return u.sessionRepo.CreateSession(ctx, userID, session)
}

func (u *SessionUseCase) MergeSessions(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*spec.MergeSessionsResponse, error) {
	cartSessionsCount, err := u.cartRepo.CountCartSessions(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	packets := make([]*models.Packet, 0)
	matchedFiltersSet := make(map[uuid.UUID]*models.SessionFilter, 0)
	sessionsEnt := make([]*ent.Session, 0)
	for paginationIndex := 0; paginationIndex < cartSessionsCount/100+1; paginationIndex++ {
		paginatedSessionsEnt, err := u.cartRepo.GetFullCart(ctx, workspaceID, paginationIndex)
		if err != nil {
			return nil, err
		}
		sessionsEnt = append(sessionsEnt, paginatedSessionsEnt...)
		for _, sessionEnt := range paginatedSessionsEnt {
			packets = append(packets, packet.PresentPacketsEnt(sessionEnt.Edges.Packets)...)
			for _, sessionFilterEnt := range sessionEnt.Edges.SessionFilters {
				filterEnt := sessionFilterEnt.Edges.Filter

				_, ok := matchedFiltersSet[filterEnt.ID]
				if !ok {
					matchedFiltersSet[filterEnt.ID] = &models.SessionFilter{
						Filter: models.Filter{
							ID:           filterEnt.ID,
							Color:        filterEnt.Color,
							InRequest:    filterEnt.InRequest,
							InResponse:   filterEnt.InResponse,
							IsBlocking:   filterEnt.IsBlocking,
							Name:         filterEnt.Name,
							Regex:        &filterEnt.Regex,
							ServiceID:    nil,
							TotalPackets: &filterEnt.TotalPackets,
							TTL:          &filterEnt.TTL,
						},
						MatchesCount: sessionFilterEnt.MatchesCount,
						SessionID:    sessionFilterEnt.SessionID,
					}
				}
			}
		}
	}

	matchedFilters := make([]*models.SessionFilter, 0, len(matchedFiltersSet))
	for _, sessionFilter := range matchedFiltersSet {
		matchedFilters = append(matchedFilters, sessionFilter)
	}

	sessionID := uuid.New()
	sessionModel := &models.Session{
		ID:                  sessionID,
		SentAt:              sessionsEnt[0].SentAt,
		ClientHost:          "pidaras_hecker",
		ServerPort:          uint16(sessionsEnt[0].ServerPort),
		Protocol:            sessionsEnt[0].Protocol,
		TTL:                 sessionsEnt[0].TTL,
		PacketsCount:        len(packets),
		Packets:             packets,
		MatchedFilters:      matchedFilters,
		ClientUserAgent:     &sessionsEnt[0].ClientUserAgent,
		AverageResponseTime: sessionsEnt[0].AverageResponseTime,
		RequestsNumber:      sessionsEnt[0].RequestsNumber,
	}
	if err := u.sessionRepo.CreateSession(ctx, userID, sessionModel); err != nil {
		return nil, err
	}
	if err := u.cartRepo.DeleteAllCartSessions(ctx, workspaceID); err != nil {
		return nil, err
	}

	return &spec.MergeSessionsResponse{
		SessionID: sessionID.String(),
	}, nil
}

func (u *SessionUseCase) GetPaginatedPcapSessions(ctx context.Context, pcapWorkspaceID uuid.UUID, paginationIndex int) (*spec.GetSessionsResponse, error) {
	sessionsEnt, err := u.sessionRepo.GetPcapSessions(ctx, pcapWorkspaceID, paginationIndex)
	if err != nil {
		return nil, err
	}
	return &spec.GetSessionsResponse{
		Sessions: session.PresentEmptySessions(sessionsEnt),
	}, nil
}

func (u *SessionUseCase) AddPcapSession(ctx context.Context, pcapWorkspaceID uuid.UUID, request *spec.AddPcapSessionRequest) error {
	session, err := session.PrepareAddSessionToPcapWorkspaceRequest(request)
	if err != nil {
		return err
	}
	return u.sessionRepo.CreatePcapSession(ctx, pcapWorkspaceID, session)
}
