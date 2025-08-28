package repositories

import (
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/pcapworkspace"
	"altron/core/repositories/ent/user"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.PcapWorkspaceRepository = (*PcapWorkspaceRepository)(nil)

type PcapWorkspaceRepository struct {
	Client *ent.Client
}

func NewPcapWorkspaceRepository(client *ent.Client) *PcapWorkspaceRepository {
	return &PcapWorkspaceRepository{
		Client: client,
	}
}

func (r *PcapWorkspaceRepository) CreatePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspace *models.PcapWorkspace) (*ent.PcapWorkspace, error) {
	return r.Client.PcapWorkspace.
		Create().
		SetUserID(userID).
		SetFileName(pcapWorkspace.FileName).
		SetStatus(pcapworkspace.Status(pcapWorkspace.Status)).
		Save(ctx)
}

func (r *PcapWorkspaceRepository) UpdatePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID, pcapWorkspace *models.PcapWorkspace) error {
	return r.Client.PcapWorkspace.Update().
		SetStatus(pcapworkspace.Status(pcapWorkspace.Status)).
		Where(pcapworkspace.And(
			pcapworkspace.IDEQ(pcapWorkspace.ID),
			pcapworkspace.HasUserWith(user.IDEQ(userID)),
		)).
		Exec(ctx)
}

func (r *PcapWorkspaceRepository) GetAllPcapWorkspaces(ctx context.Context, userID uuid.UUID) ([]*ent.PcapWorkspace, error) {
	return r.Client.PcapWorkspace.
		Query().
		Where(pcapworkspace.HasUserWith(user.IDEQ(userID))).
		All(ctx)
}

func (r *PcapWorkspaceRepository) GetPcapWorkspaceByID(ctx context.Context, userID, pcapWorkspaceID uuid.UUID) (*ent.PcapWorkspace, error) {
	pcapWorkspace, err := r.Client.PcapWorkspace.
		Query().
		Where(pcapworkspace.And(
			pcapworkspace.IDEQ(pcapWorkspaceID),
			pcapworkspace.HasUserWith(user.IDEQ(userID)),
		)).Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, models.ErrorPcapWorkspaceNotFound
		}
		return nil, err
	}
	return pcapWorkspace, nil
}

func (r *PcapWorkspaceRepository) DeletePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID) error {
	_, err := r.Client.PcapWorkspace.
		Delete().
		Where(pcapworkspace.And(
			pcapworkspace.IDEQ(pcapWorkspaceID),
			pcapworkspace.HasUserWith(user.IDEQ(userID)),
		)).
		Exec(ctx)
	return err
}
