package usecases

import (
	common "altron/common/models"
	"altron/config"
	"altron/core/adapter/workspace"
	"altron/core/dto"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"context"
	"fmt"
	"strings"

	"altron/pkg/redis"
	req "altron/pkg/request"

	"github.com/google/uuid"
)

var _ interfaces.WorkspaceUseCase = (*WorkspaceUseCase)(nil)

type WorkspaceUseCase struct {
	cfg           *config.AppConfig
	redisClient   *redis.Client[[]spec.Characteristic]
	workspaceRepo interfaces.WorkspaceRepository
	sessionRepo   interfaces.SessionRepository
}

func NewWorkspaceUseCase(
	cfg *config.AppConfig,
	redisClient *redis.Client[[]spec.Characteristic],
	workspaceRepo interfaces.WorkspaceRepository,
	sessionRepo interfaces.SessionRepository,
) *WorkspaceUseCase {
	return &WorkspaceUseCase{
		cfg:           cfg,
		redisClient:   redisClient,
		workspaceRepo: workspaceRepo,
		sessionRepo:   sessionRepo,
	}
}

func (u *WorkspaceUseCase) CreateWorkspace(ctx context.Context, userID uuid.UUID, request *spec.CreateWorkspaceRequest) (*spec.CreateWorkspaceResponse, error) {
	workspaceEnt, err := u.workspaceRepo.CreateWorkspace(ctx, userID, workspace.PrepareCreateWorkspaceRequest(request))
	if err != nil {
		return nil, err
	}
	serviceEnt := workspaceEnt.Edges.Service
	var workspaceID string
	if strings.EqualFold(workspaceEnt.Name, "checker") {
		workspaceID = "checker"
	} else {
		workspaceID = workspaceEnt.ID.String()
	}
	for _, componentName := range common.ComponentNames {
		if err := u.redisClient.Set(
			ctx,
			fmt.Sprintf("%d_%s_%s", serviceEnt.Port, workspaceID, componentName),
			[]spec.Characteristic{},
		); err != nil {
			return nil, err
		}
	}
	if err := req.Put(
		fmt.Sprintf("http://%s:%d/workspaces/%s/sessions/%d",
			u.cfg.AltronHost,
			u.cfg.AltronConnectionPort,
			workspaceEnt.ID,
			request.ServicePort),
		dto.StartSessionsListeningsRequest{
			Timeout:   request.Timeout,
			StartTime: request.StartTime,
			IsChecker: strings.EqualFold(workspaceEnt.Name, "checker"),
		},
	); err != nil {
		return nil, err
	}
	return &spec.CreateWorkspaceResponse{
		Id:          workspaceEnt.ID.String(),
		Name:        workspaceEnt.Name,
		ServicePort: request.ServicePort,
	}, nil
}

func (u *WorkspaceUseCase) UpdateWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, request *spec.UpdateWorkspaceRequest) error {
	return u.workspaceRepo.UpdateWorkspace(
		ctx,
		userID,
		workspaceID,
		workspace.PrepareUpdateWorkspaceRequest(request),
	)
}

func (u *WorkspaceUseCase) DeleteWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) error {
	workspaceEnt, err := u.workspaceRepo.DeleteWorkspace(ctx, userID, workspaceID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(string(workspaceEnt.Status), "COMPLETED") {
		_, err := req.DeleteWithEmptyResponse(
			fmt.Sprintf("http://%s:%d/workspaces/%s/sessions", u.cfg.AltronHost, u.cfg.AltronConnectionPort, workspaceID),
		)
		return err
	}
	serviceEnt := workspaceEnt.Edges.Service
	var workspaceId string
	if strings.EqualFold(workspaceEnt.Name, "checker") {
		workspaceId = "checker"
	} else {
		workspaceId = workspaceID.String()
	}
	for _, componentName := range common.ComponentNames {
		if err := u.redisClient.Set(
			ctx,
			fmt.Sprintf("%d_%s_%s", serviceEnt.Port, workspaceId, componentName),
			[]spec.Characteristic{},
		); err != nil {
			return err
		}
	}
	return nil
}

func (u *WorkspaceUseCase) GetWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*spec.GetWorkspaceResponse, error) {
	workspace, err := u.workspaceRepo.GetWorkspaceByID(ctx, userID, workspaceID)
	if err != nil {
		return nil, err
	}

	return &spec.GetWorkspaceResponse{
		Workspace: spec.Workspace{
			Id:          workspace.ID.String(),
			Name:        workspace.Name,
			ServicePort: uint16(workspace.Edges.Service.Port),
			Status:      spec.WorkspaceStatus(workspace.Status),
		},
	}, nil
}

func (u *WorkspaceUseCase) ResetWorkspaces(ctx context.Context, userID uuid.UUID) error {
	return u.workspaceRepo.UpdateAllWorkspaces(ctx, userID, common.COMPLETED)
}
