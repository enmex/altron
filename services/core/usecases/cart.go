package usecases

import (
	"altron/core/adapter/session"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"context"

	"github.com/google/uuid"
)

var _ interfaces.CartUseCase = (*CartUseCase)(nil)

type CartUseCase struct {
	cartRepo interfaces.CartRepository
}

func NewCartUseCase(cartRepo interfaces.CartRepository) *CartUseCase {
	return &CartUseCase{
		cartRepo: cartRepo,
	}
}

func (u *CartUseCase) AddSessions(ctx context.Context, workspaceID uuid.UUID, request *spec.AddSessionsToCartRequest) error {
	sessionIDs := make([]uuid.UUID, 0, len(request.Sessions))
	for _, session := range request.Sessions {
		sessionIDs = append(sessionIDs, uuid.MustParse(session))
	}
	return u.cartRepo.AddSessions(ctx, workspaceID, sessionIDs)
}

func (u *CartUseCase) GetCart(ctx context.Context, workspaceID uuid.UUID, paginationIndex int) (*spec.GetCartResponse, error) {
	sessionsEnt, err := u.cartRepo.GetCart(ctx, workspaceID, paginationIndex)
	if err != nil {
		return nil, err
	}
	return &spec.GetCartResponse{
		Sessions: session.PresentEmptySessions(sessionsEnt),
	}, nil
}

func (u *CartUseCase) RemoveSession(ctx context.Context, workspaceID, sessionID uuid.UUID) error {
	return u.cartRepo.DeleteSession(ctx, workspaceID, sessionID)
}

func (u *CartUseCase) ClearCart(ctx context.Context, workspaceID uuid.UUID) error {
	return u.cartRepo.DeleteAllCartSessions(ctx, workspaceID)
}
