package repositories

import (
	common "altron/common/models"
	"altron/core/interfaces"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/analyzercomponent"
	"altron/core/repositories/ent/analyzerpayload"
	"altron/core/repositories/ent/service"
	"altron/core/repositories/ent/user"
	"altron/core/repositories/ent/workspace"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.AnalyzerPayloadRepository = (*AnalyzerPayloadRepository)(nil)

type AnalyzerPayloadRepository struct {
	Client *ent.Client
}

func NewAnalyzerPayloadRepository(client *ent.Client) *AnalyzerPayloadRepository {
	return &AnalyzerPayloadRepository{
		Client: client,
	}
}

func (r *AnalyzerPayloadRepository) CreateAnalyzerPayloads(ctx context.Context, workspaceID uuid.UUID, analyzerPayload common.AnalyzerPayload) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}

	for componentName, characteristics := range analyzerPayload.Payload {
		analyzerComponentEnt, err := tx.AnalyzerComponent.
			Query().
			Where(analyzercomponent.NameEQ(componentName)).
			Only(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}

		for _, characteristic := range characteristics {
			err := tx.AnalyzerPayload.Create().
				SetAnalyzerComponent(analyzerComponentEnt).
				SetWorkspaceID(workspaceID).
				SetValue(characteristic.Value).
				SetNumber(characteristic.Number).
				Exec(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *AnalyzerPayloadRepository) GetAnalyzerPayloads(ctx context.Context, userID uuid.UUID, servicePort uint16, workspaceID uuid.UUID, componentName string) ([]*ent.AnalyzerPayload, error) {
	payload, err := r.Client.AnalyzerPayload.Query().
		Where(analyzerpayload.And(
			analyzerpayload.HasAnalyzerComponentWith(analyzercomponent.NameEQ(componentName)),
			analyzerpayload.HasWorkspaceWith(workspace.And(
				workspace.IDEQ(workspaceID),
				workspace.HasServiceWith(service.And(
					service.PortEQ(uint32(servicePort)),
					service.HasUserWith(user.IDEQ(userID)),
				)),
			)),
		)).
		All(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return make([]*ent.AnalyzerPayload, 0), nil
		}
		return nil, err
	}

	return payload, nil
}

func (r *AnalyzerPayloadRepository) GetAnalyzerComponents(ctx context.Context) ([]*ent.AnalyzerComponent, error) {
	return r.Client.AnalyzerComponent.Query().All(ctx)
}

func (r *AnalyzerPayloadRepository) CreateAnalyzerComponents(ctx context.Context, components []string) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}

	for _, characteristic := range components {
		_, err := tx.AnalyzerComponent.
			Query().
			Where(analyzercomponent.NameEQ(characteristic)).
			Only(ctx)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
			err := tx.AnalyzerComponent.
				Create().
				SetName(characteristic).
				Exec(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	return tx.Commit()
}
