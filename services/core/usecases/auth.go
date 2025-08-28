package usecases

import (
	"altron/config"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/pkg/auth"
	"altron/utils"
	"context"

	common "altron/common/models"
)

var _ interfaces.AuthUseCase = (*AuthUseCase)(nil)

type AuthUseCase struct {
	cfg     *config.Config
	jwtAuth *auth.JwtAuthenticate
	userRepo    interfaces.UserRepository
	workspaceRepo interfaces.WorkspaceRepository
}

func NewAuthUseCase(
	cfg *config.Config,
	jwtAuth *auth.JwtAuthenticate,
	userRepo interfaces.UserRepository,
	workspaceRepo interfaces.WorkspaceRepository,
) *AuthUseCase {
	return &AuthUseCase{
		cfg:     cfg,
		jwtAuth: jwtAuth,
		userRepo:    userRepo,
		workspaceRepo: workspaceRepo,
	}
}

func (u *AuthUseCase) SignUp(ctx context.Context, payload *spec.AuthRequest) (*spec.AuthResponse, error) {
	enc, err := utils.CryptString(payload.Password, u.cfg.Auth.HashSalt)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: u.cfg.App.AltronName,
		Password: enc,
	}
	userEnt, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	if _, err := u.workspaceRepo.CreateUserWorkspace(ctx, &models.Workspace{
		Name: userEnt.ID.String(),
		Status: common.COMPLETED,
	}); err != nil {
		return nil, err
	}

	claims := auth.AccessClaims{
		UserID: userEnt.ID,
	}

	jwtToken, err := u.jwtAuth.GenerateAccessToken(claims)

	if err != nil {
		return nil, err
	}

	return &spec.AuthResponse{
		Id:    userEnt.ID.String(),
		Token: jwtToken,
	}, nil
}

func (u *AuthUseCase) SignIn(ctx context.Context, payload *spec.AuthRequest) (*spec.AuthResponse, error) {
	user, err := u.userRepo.GetUser(ctx, u.cfg.App.AltronName)
	if err != nil {
		return nil, models.ErrorInvalidCredentials
	}

	dec, err := utils.DecryptString(user.Password, u.cfg.Auth.HashSalt)
	if err != nil {
		return nil, err
	}
	if dec != payload.Password {
		return nil, models.ErrorInvalidCredentials
	}

	claims := auth.AccessClaims{
		UserID: user.ID,
	}

	jwtToken, err := u.jwtAuth.GenerateAccessToken(claims)

	if err != nil {
		return nil, err
	}

	return &spec.AuthResponse{
		Id:    user.ID.String(),
		Token: jwtToken,
	}, nil
}

func (u *AuthUseCase) AuthUser(ctx context.Context) (*spec.AuthResponse, error) {
	user, err := u.userRepo.GetUser(ctx, u.cfg.App.AltronName)
	if err != nil {
		return nil, models.ErrorInvalidCredentials
	}

	claims := auth.AccessClaims{
		UserID: user.ID,
	}

	jwtToken, err := u.jwtAuth.GenerateAccessToken(claims)

	if err != nil {
		return nil, err
	}

	return &spec.AuthResponse{
		Id:    user.ID.String(),
		Token: jwtToken,
	}, nil
}

func (u *AuthUseCase) SignUpManager(ctx context.Context) (*spec.AuthResponse, error) {
	enc, err := utils.CryptString(u.cfg.App.ManagerPassword, u.cfg.Auth.HashSalt)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: "manager",
		Password: enc,
	}
	userEnt, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	claims := auth.AccessClaims{
		UserID: userEnt.ID,
	}

	jwtToken, err := u.jwtAuth.GenerateAccessToken(claims)

	if err != nil {
		return nil, err
	}

	return &spec.AuthResponse{
		Id:    userEnt.ID.String(),
		Token: jwtToken,
	}, nil
}

func (u *AuthUseCase) SignInManager(ctx context.Context, payload *spec.AuthRequest) (*spec.AuthResponse, error) {
	user, err := u.userRepo.GetUser(ctx, "manager")
	if err != nil {
		return nil, models.ErrorInvalidCredentials
	}
	
	dec, err := utils.DecryptString(user.Password, u.cfg.Auth.HashSalt)
	if err != nil {
		return nil, err
	}
	if dec != payload.Password {
		return nil, models.ErrorInvalidCredentials
	}

	claims := auth.AccessClaims{
		UserID: user.ID,
	}

	jwtToken, err := u.jwtAuth.GenerateAccessToken(claims)

	if err != nil {
		return nil, err
	}

	return &spec.AuthResponse{
		Id:    user.ID.String(),
		Token: jwtToken,
	}, nil
}

func (u *AuthUseCase) GetUserInfo(ctx context.Context) (*spec.GetUserInfoResponse, error) {
	user, err := u.userRepo.GetUser(ctx, u.cfg.App.AltronName)
	if err != nil {
		return nil, err
	}

	return &spec.GetUserInfoResponse{
		Username: user.Username,
		Password: user.Password,
	}, nil
}
