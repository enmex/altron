package repositories

import (
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/user"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	Client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		Client: client,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*ent.User, error) {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	userEnt, err := tx.User.
		Create().
		SetUsername(user.Username).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		if strings.Contains(err.Error(), "duplicate") {
			return nil, models.ErrorUserAlreadyExists
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return userEnt, nil
}

func (r *UserRepository) GetUser(ctx context.Context, username string) (*ent.User, error) {
	userEnt, err := r.Client.User.
		Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	return userEnt, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userModel *models.User) error {
	_, err := r.Client.User.
		Update().
		SetPassword(userModel.Password).
		Where(user.UsernameEQ(userModel.Username)).
		Save(ctx)
	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	user, err := r.Client.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, models.ErrorUserNotFound
		}
		return nil, err
	}
	return user, nil
}
