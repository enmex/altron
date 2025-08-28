package usecases

import (
	common "altron/common/models"
	"altron/core/adapter/filter"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"context"
	"encoding/base64"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var _ interfaces.FilterUseCase = (*FilterUseCase)(nil)

type FilterUseCase struct {
	log         *logrus.Logger
	filterRepo  interfaces.FilterRepository
	sessionRepo interfaces.SessionRepository
}

func NewFilterUseCase(
	log *logrus.Logger,
	filterRepo interfaces.FilterRepository,
	sessionRepo interfaces.SessionRepository,
) *FilterUseCase {
	return &FilterUseCase{
		log:         log,
		filterRepo:  filterRepo,
		sessionRepo: sessionRepo,
	}
}

func (u *FilterUseCase) CreateFilter(ctx context.Context, userID uuid.UUID, request *spec.CreateFilterRequest) error {
	var serviceUuid *uuid.UUID
	if request.ServiceId != nil {
		uuid := uuid.MustParse(*request.ServiceId)
		serviceUuid = &uuid
	}

	filter := &common.Filter{
		Name:         request.Name,
		ServiceID:    serviceUuid,
		Regex:        request.Regex,
		TTL:          request.Ttl,
		Color:        request.Color,
		TotalPackets: request.TotalPackets,
		InRequest:    request.InRequest,
		InResponse:   request.InResponse,
		IsBlocking:   request.IsBlocking,
	}
	filterEnt, err := u.filterRepo.CreateFilter(ctx, userID, filter)
	if err != nil {
		return err
	}
	if filter.IsBlocking {
		return nil
	}
	filter.ID = filterEnt.ID

	sessionFilters, err := u.createSessionFilters(ctx, userID, filter)
	if err != nil {
		return err
	}

	return u.filterRepo.CreateSessionFilters(ctx, sessionFilters)
}

func (u *FilterUseCase) GetAllFilters(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetAllFiltersResponse, error) {
	filtersEnt, err := u.filterRepo.GetAllFilters(ctx, userID, servicePort)
	if err != nil {
		return nil, err
	}

	return &spec.GetAllFiltersResponse{
		Filters: filter.PresentFilters(filtersEnt),
	}, nil
}

func (u *FilterUseCase) DeleteFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID) error {
	_, err := u.filterRepo.DeleteFilter(ctx, userID, filterID)
	return err
}

func (u *FilterUseCase) UpdateFilter(ctx context.Context, userID uuid.UUID, request *spec.UpdateFilterRequest, filterID uuid.UUID) error {
	if err := u.filterRepo.DeleteSessionFilter(ctx, userID, filterID); err != nil {
		return err
	}
	filter := filter.PrepareUpdateFilterRequest(request, filterID)
	if err := u.filterRepo.UpdateFilter(ctx, userID, filter); err != nil {
		return err
	}

	sessionFilters, err := u.createSessionFilters(ctx, userID, filter)
	if err != nil {
		return err
	}
	return u.filterRepo.CreateSessionFilters(ctx, sessionFilters)
}

func (u *FilterUseCase) createSessionFilters(ctx context.Context, userID uuid.UUID, filter *common.Filter) ([]*common.SessionFilter, error) {
	sessionsCount, err := u.sessionRepo.CountSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sessionsCount == 0 {
		return nil, nil
	}
	sessionFilters := make([]*common.SessionFilter, 0)
	paginationCount := 1 + sessionsCount/100
	wg := &sync.WaitGroup{}
	wg.Add(paginationCount)

	for paginationIndex := 0; paginationIndex < paginationCount; paginationIndex++ {
		go func(paginationIndex int) {
			defer wg.Done()
			sessions, err := u.sessionRepo.GetSessions(ctx, userID, paginationIndex)
			if err != nil {
				u.log.Errorln(err)
				return
			}
			for _, session := range sessions {
				matchesCount := 0
				matchedPackets := make([]int, 0)

				match := true
				if filter.TTL != nil {
					match = session.TTL == *filter.TTL
				}
				if filter.TotalPackets != nil {
					match = match && session.PacketsCount/2 == *filter.TotalPackets
				}
				if filter.Regex != nil {
					re := regexp.MustCompile(*filter.Regex)

					for idx, packet := range session.Edges.Packets {
						payload, err := base64.StdEncoding.DecodeString(packet.Payload)
						if err != nil {
							u.log.Errorln(err)
							break
						}
						matches := re.FindAllStringIndex(string(payload), -1)
						if !match {
							match = len(matches) != 0
						}
						if len(matches) > 0 {
							matchesCount += len(matches)
							matchedPackets = append(matchedPackets, idx)
						}
					}
				}
				if match {
					sessionFilters = append(sessionFilters, &common.SessionFilter{
						Filter:         *filter,
						MatchesCount:   matchesCount,
						MatchedPackets: matchedPackets,
						SessionID:      session.ID,
					})
				}
			}
		}(paginationIndex)
	}
	wg.Wait()

	return sessionFilters, nil
}
