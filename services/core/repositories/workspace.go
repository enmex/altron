package repositories

import (
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/service"
	"altron/core/repositories/ent/user"
	"altron/core/repositories/ent/workspace"
	"context"
	"strings"

	common "altron/common/models"
	"github.com/google/uuid"
)

var _ interfaces.WorkspaceRepository = (*WorkspaceRepository)(nil)

type WorkspaceRepository struct {
	Client *ent.Client
}

func NewWorkspaceRepository(client *ent.Client) *WorkspaceRepository {
	return &WorkspaceRepository{
		Client: client,
	}
}

func (r *WorkspaceRepository) CreateUserWorkspace(ctx context.Context, workspaceModel *models.Workspace) (*ent.Workspace, error) {
	return r.Client.Workspace.Create().
		SetName(workspaceModel.Name).
		SetStatus(workspace.Status(workspaceModel.Status)).
		Save(ctx)
}

func (r *WorkspaceRepository) CreateWorkspace(ctx context.Context, userID uuid.UUID, workspaceModel *models.Workspace) (*ent.Workspace, error) {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	serviceEnt, err := tx.Service.
		Query().
		Where(service.And(
			service.PortEQ(uint32(workspaceModel.ServicePort)),
			service.HasUserWith(user.IDEQ(userID)),
		)).
		Only(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = tx.Workspace.
		Create().
		SetName(workspaceModel.Name).
		SetStatus(workspace.Status(workspaceModel.Status)).
		SetService(serviceEnt).
		Save(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		if strings.Contains(err.Error(), "duplicate") {
			return nil, models.ErrorWorkspaceAlreadyExists
		}
		return nil, err
	}

	workspaceEnt, err := tx.Workspace.
		Query().
		WithService().
		Where(workspace.And(
			workspace.NameEQ(workspaceModel.Name),
			workspace.HasServiceWith(service.IDEQ(serviceEnt.ID))),
		).
		Only(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return workspaceEnt, nil
}

func (r *WorkspaceRepository) UpdateWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID, workspaceModel *models.Workspace) error {
	qb := r.Client.Workspace.
		Update()
	if len(workspaceModel.Name) != 0 {
		qb = qb.SetName(workspaceModel.Name)
	}
	if len(workspaceModel.Status) != 0 {
		qb = qb.SetStatus(workspace.Status(workspaceModel.Status))
	}
	err := qb.SetStatus(workspace.Status(workspaceModel.Status)).
		Where(workspace.And(
			workspace.IDEQ(workspaceID),
			workspace.HasServiceWith(service.HasUserWith(user.IDEQ(userID))),
		)).
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.ErrorWorkspaceNotFound
		}
		return err
	}
	return nil
}

func (r *WorkspaceRepository) DeleteWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*ent.Workspace, error) {
	workspaceEnt, err := r.Client.Workspace.
		Query().
		WithService().
		Where(workspace.And(
			workspace.IDEQ(workspaceID),
			workspace.HasServiceWith(service.HasUserWith(user.IDEQ(userID))),
		)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.Client.Workspace.
		Delete().
		Where(workspace.IDEQ(workspaceID)).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return workspaceEnt, nil
}

func (r *WorkspaceRepository) GetWorkspacesByServicePort(ctx context.Context, userID uuid.UUID, servicePort uint16) ([]*ent.Workspace, error) {
	return r.Client.Workspace.
		Query().
		Where(workspace.HasServiceWith(service.And(
			service.PortEQ(uint32(servicePort)),
			service.HasUserWith(user.IDEQ(userID)),
		))).
		All(ctx)
}

func (r *WorkspaceRepository) GetServiceWorkspaceByName(ctx context.Context, userID uuid.UUID, servicePort uint16, workspaceName string) (*ent.Workspace, error) {
	workspaceEnt, err := r.Client.Workspace.
		Query().
		WithService().
		Where(workspace.And(
			workspace.NameEQ(workspaceName),
			workspace.HasServiceWith(service.And(
				service.PortEQ(uint32(servicePort)),
				service.HasUserWith(user.IDEQ(userID)),
			)),
		)).
		Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, models.ErrorWorkspaceNotFound
		}
		return nil, err
	}

	return workspaceEnt, nil
}

func (r *WorkspaceRepository) GetWorkspaceByID(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (*ent.Workspace, error) {
	workspaceEnt, err := r.Client.Workspace.
		Query().
		WithService().
		Where(workspace.And(
			workspace.IDEQ(workspaceID),
			workspace.HasServiceWith(service.HasUserWith(user.IDEQ(userID))),
		)).
		Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, models.ErrorWorkspaceNotFound
		}
		return nil, err
	}

	return workspaceEnt, nil
}

func (r *WorkspaceRepository) UpdateAllWorkspaces(ctx context.Context, userID uuid.UUID, status common.WorkspaceStatus) error {
	_, err := r.Client.Workspace.
		Update().SetStatus(workspace.Status(status)).
		Where(workspace.HasServiceWith(service.HasUserWith(user.IDEQ(userID)))).
		Save(ctx)
	return err
}
