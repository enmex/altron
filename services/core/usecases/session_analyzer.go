package usecases

import (
	common "altron/common/models"
	analyzerpayload "altron/core/adapter/analyzer_payload"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/pkg/redis"
	"altron/utils"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var _ interfaces.SessionAnalyzerUseCase = (*SessionAnalyzerUseCase)(nil)

type SessionAnalyzerUseCase struct {
	analyzerPayloadRepo interfaces.AnalyzerPayloadRepository
	workspaceRepo       interfaces.WorkspaceRepository
	redisClient         *redis.Client[[]spec.Characteristic]
}

func NewSessionAnalyzerUseCase(
	analyzerPayloadRepo interfaces.AnalyzerPayloadRepository,
	workspaceRepo interfaces.WorkspaceRepository,
	redisClient *redis.Client[[]spec.Characteristic],
) *SessionAnalyzerUseCase {
	return &SessionAnalyzerUseCase{
		analyzerPayloadRepo: analyzerPayloadRepo,
		workspaceRepo:       workspaceRepo,
		redisClient:         redisClient,
	}
}

func (u *SessionAnalyzerUseCase) CreateAnalyzerPayload(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, request *spec.CreateAnalyzerPayloadRequest) error {
	return u.analyzerPayloadRepo.CreateAnalyzerPayloads(ctx, workspaceID, *analyzerpayload.PrepareCreateAnalyzerPayloadRequest(request))
}

func (u *SessionAnalyzerUseCase) GetAllAnalyzerComponents(ctx context.Context) (*spec.GetAnalyzerComponentsResponse, error) {
	componentsEnt, err := u.analyzerPayloadRepo.GetAnalyzerComponents(ctx)
	if err != nil {
		return nil, err
	}

	componentNames := make([]string, 0, len(componentsEnt))
	for _, componentEnt := range componentsEnt {
		componentNames = append(componentNames, componentEnt.Name)
	}

	return &spec.GetAnalyzerComponentsResponse{
		Components: componentNames,
	}, nil
}

func (u *SessionAnalyzerUseCase) GetAnalyzerPayload(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetSessionAnalyzerResponse, error) {
	payload := make(map[string][]spec.Characteristic)

	checkerMask := make(map[string][]*ent.AnalyzerPayload)
	workspace, err := u.workspaceRepo.GetServiceWorkspaceByName(ctx, userID, servicePort, "checker")
	hasChecker := err == nil
	if hasChecker {
		for _, component := range common.ComponentNames {
			componentCharacteristics, err := u.analyzerPayloadRepo.GetAnalyzerPayloads(ctx, userID, servicePort, workspace.ID, component)
			if err != nil {
				return nil, err
			}

			checkerMask[component] = componentCharacteristics
		}
	}
	for _, component := range common.ComponentNames {
		data, err := u.redisClient.Get(ctx, fmt.Sprintf("%d_%s", servicePort, component))
		if err != nil {
			return nil, err
		}
		componentCharacteristics := *data
		for idx := range componentCharacteristics {
			checkerComponentPayload, ok := checkerMask[component]
			componentCharacteristics[idx].IsSafe = ok && utils.ContainsFunc[*ent.AnalyzerPayload](
				checkerComponentPayload,
				func(c *ent.AnalyzerPayload) bool {
					return c.Value == componentCharacteristics[idx].Value
				},
			)
		}
		payload[component] = componentCharacteristics
	}
	return &spec.GetSessionAnalyzerResponse{
		HasChecker: hasChecker,
		Analyzer: spec.GetSessionAnalyzerResponse_Analyzer{
			AdditionalProperties: payload,
		},
	}, nil
}

func (u *SessionAnalyzerUseCase) GetWorkspaceAnalyzerPayload(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*spec.GetSessionAnalyzerResponse, error) {
	payload := make(map[string][]spec.Characteristic)

	workspace, err := u.workspaceRepo.GetWorkspaceByID(ctx, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	isChecker := workspace.Name == "checker"
	for _, component := range common.ComponentNames {
		componentPayload := make([]spec.Characteristic, 0)
		if workspace.Status == "LISTENING" {
			var workspaceID string
			if isChecker {
				workspaceID = "checker"
			} else {
				workspaceID = workspace.ID.String()
			}
			data, err := u.redisClient.Get(ctx, fmt.Sprintf("%d_%s_%s", workspace.Edges.Service.Port, workspaceID, component))
			if err != nil {
				return nil, err
			}
			componentPayload = *data
		} else {
			componentPayloadEnt, err := u.analyzerPayloadRepo.GetAnalyzerPayloads(ctx, userID, uint16(workspace.Edges.Service.Port), workspaceID, component)
			if err != nil {
				return nil, err
			}
			for _, ch := range componentPayloadEnt {
				componentPayload = append(componentPayload, spec.Characteristic{
					Value:  ch.Value,
					Number: ch.Number,
					IsSafe: isChecker,
				})
			}
		}
		payload[component] = componentPayload
	}
	return &spec.GetSessionAnalyzerResponse{
		HasChecker: isChecker,
		Analyzer: spec.GetSessionAnalyzerResponse_Analyzer{
			AdditionalProperties: payload,
		},
	}, nil
}

func (u *SessionAnalyzerUseCase) GetServiceCheckerMask(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetSessionAnalyzerResponse, error) {
	payload := make(map[string][]spec.Characteristic)

	workspace, err := u.workspaceRepo.GetServiceWorkspaceByName(ctx, userID, servicePort, "checker")
	if err != nil && errors.Is(err, models.ErrorWorkspaceNotFound) {
		return &spec.GetSessionAnalyzerResponse{
			HasChecker: false,
			Analyzer: spec.GetSessionAnalyzerResponse_Analyzer{
				AdditionalProperties: payload,
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	for _, component := range common.ComponentNames {
		componentPayload, err := u.analyzerPayloadRepo.GetAnalyzerPayloads(ctx, userID, servicePort, workspace.ID, component)
		if err != nil {
			return nil, err
		}
		chs := make([]spec.Characteristic, 0, len(componentPayload))
		for _, ch := range componentPayload {
			chs = append(chs, spec.Characteristic{
				Value:  ch.Value,
				Number: ch.Number,
				IsSafe: true,
			})
		}
		payload[component] = chs
	}
	return &spec.GetSessionAnalyzerResponse{
		HasChecker: true,
		Analyzer: spec.GetSessionAnalyzerResponse_Analyzer{
			AdditionalProperties: payload,
		},
	}, nil
}
