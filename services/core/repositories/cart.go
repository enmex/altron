package repositories

import (
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/cartsession"
	"altron/core/repositories/ent/workspace"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.CartRepository = (*CartRepository)(nil)

type CartRepository struct {
	Client *ent.Client
}

func NewCartRepository(client *ent.Client) *CartRepository {
	return &CartRepository{
		Client: client,
	}
}

func (r *CartRepository) AddSessions(ctx context.Context, workspaceID uuid.UUID, sessionIDs []uuid.UUID) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}
	for _, sessionID := range sessionIDs {
		_, err := tx.CartSession.Create().
			SetSessionID(sessionID).
			SetWorkspaceID(workspaceID).
			Save(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			if strings.Contains(err.Error(), "duplicate") {
				return models.ErrorDuplicateSession
			}
			return err
		}
	}
	return tx.Commit()
}

func (r *CartRepository) DeleteSession(ctx context.Context, workspaceID uuid.UUID, sessionID uuid.UUID) error {
	return r.Client.Workspace.UpdateOneID(workspaceID).
		RemoveCartSessionIDs(sessionID).
		Exec(ctx)
}

func (r *CartRepository) GetCart(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) ([]*ent.Session, error) {
	workspace, err := r.Client.Workspace.Query().
		Where(workspace.IDEQ(workspaceID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	cartSessions, err := r.Client.Workspace.QueryCartSessions(workspace).
		QuerySession().
		WithSessionFilters(func(sfq *ent.SessionFilterQuery) {
			sfq.WithFilter()
		}).
		Order(ent.Asc("sent_at")).
		Offset(paginationIndex * 100).
		Limit(100).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return cartSessions, nil
}

func (r *CartRepository) DeleteAllCartSessions(ctx context.Context, workspaceID uuid.UUID) error {
	_, err := r.Client.CartSession.Delete().
		Where(cartsession.HasWorkspaceWith(workspace.IDEQ(workspaceID))).
		Exec(ctx)
	return err
}

func (r *CartRepository) CountCartSessions(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	workspace, err := r.Client.Workspace.Query().
		Where(workspace.IDEQ(workspaceID)).
		Only(ctx)
	if err != nil {
		return 0, err
	}

	return r.Client.Workspace.QueryCartSessions(workspace).
		Count(ctx)
}

func (r *CartRepository) GetFullCart(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) ([]*ent.Session, error) {
	workspace, err := r.Client.Workspace.Query().
		Where(workspace.IDEQ(workspaceID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	cartSessions, err := r.Client.Workspace.QueryCartSessions(workspace).
		QuerySession().
		WithPackets().
		WithSessionFilters(func(sfq *ent.SessionFilterQuery) {
			sfq.WithFilter()
		}).
		Order(ent.Asc("sent_at")).
		Offset(paginationIndex * 100).
		Limit(100).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return cartSessions, nil
}
