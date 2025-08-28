package repositories

import (
	"altron/core/interfaces"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/filter"
	"altron/core/repositories/ent/pcapworkspace"
	"altron/core/repositories/ent/service"
	u "altron/core/repositories/ent/user"
	"altron/core/repositories/ent/workspace"
	"context"
)

var _ interfaces.AdminRepository = (*AdminRepository)(nil)

type AdminRepository struct {
	Client *ent.Client
}

func NewAdminRepository(client *ent.Client) *AdminRepository {
	return &AdminRepository{
		Client: client,
	}
}

func (r *AdminRepository) ClearAll(ctx context.Context, username string) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}
	user, err := r.Client.User.Query().
		Where(u.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if _, err := tx.Service.Delete().Where(service.HasUserWith(u.IDEQ(user.ID))).Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if _, err := tx.Workspace.Delete().Where(workspace.HasServiceWith(service.HasUserWith(u.IDEQ(user.ID)))).Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if _, err := tx.Filter.Delete().Where(filter.HasUserWith(u.IDEQ(user.ID))).Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if _, err := tx.PcapWorkspace.Delete().Where(pcapworkspace.HasUserWith(u.IDEQ(user.ID))).Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
