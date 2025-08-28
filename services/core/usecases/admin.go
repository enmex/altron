package usecases

import (
	"altron/config"
	"altron/core/dto"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var _ interfaces.AdminUseCase = (*AdminUseCase)(nil)

type AdminUseCase struct {
	cfg       *config.Config
	userRepo  interfaces.UserRepository
	adminRepo interfaces.AdminRepository
}

func NewAdminUseCase(
	cfg *config.Config,
	userRepo interfaces.UserRepository,
	adminRepo interfaces.AdminRepository,
) *AdminUseCase {
	return &AdminUseCase{
		cfg:       cfg,
		userRepo:  userRepo,
		adminRepo: adminRepo,
	}
}

func (u *AdminUseCase) ResetAll(ctx context.Context) (*spec.ReconfigureAltronResponse, error) {
	if err := u.adminRepo.ClearAll(ctx, u.cfg.App.AltronName); err != nil {
		return nil, err
	}
	newPassword := utils.RandomString(12)
	enc, err := utils.CryptString(newPassword, u.cfg.Auth.HashSalt)
	if err != nil {
		return nil, err
	}
	if err := u.userRepo.UpdateUser(ctx, &models.User{
		Username: u.cfg.App.AltronName,
		Password: enc,
	}); err != nil {
		return nil, err
	}
	u.cfg.Auth.Secret = []byte(utils.RandomString(20))

	client := &http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s/api/v2/token", u.cfg.SFTP.APIServer),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("admin", u.cfg.SFTP.AdminPassword)

	authRes, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if authRes.StatusCode != http.StatusOK {
		return nil, models.ErrorSomethingWrong
	}
	defer authRes.Body.Close()

	var authResponse dto.SFTPTokenResponse
	if err := json.NewDecoder(authRes.Body).Decode(&authResponse); err != nil {
		return nil, err
	}
	username := strings.ReplaceAll(u.cfg.App.AltronName, "-", "")

	data, err := json.Marshal(dto.ChangeUserPasswordRequest{
		HomeDir:  "/var/lib/sftpgo/" + username,
		Password: newPassword,
		Permissions: map[string][]string{
			"/": {"*"},
		},
		Status: 1,
	})
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("http://%s/api/v2/users/%s", u.cfg.SFTP.APIServer, username),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	updateRes, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if updateRes.StatusCode != http.StatusOK {
		return nil, models.ErrorSomethingWrong
	}

	data, err = json.Marshal(dto.CreateEventRequest{
		Type:            models.AltronReset,
		WaitForResponse: false,
	})
	if err != nil {
		return nil, err
	}
	if _, err := http.Post(
		fmt.Sprintf("http://%s:%d/events", u.cfg.App.AltronHost, u.cfg.App.AltronConnectionPort),
		"application/json",
		bytes.NewBuffer(data),
	); err != nil {
		return nil, err
	}

	return &spec.ReconfigureAltronResponse{
		Password: newPassword,
	}, nil
}

func (u *AdminUseCase) GetInfo(ctx context.Context) (*spec.GetInfoResponse, error) {
	return &spec.GetInfoResponse{
		Username: strings.ReplaceAll(u.cfg.App.AltronName, "-", ""),
	}, nil
}
