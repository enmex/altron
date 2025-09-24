package usecases

import (
	common "altron/common/models"
	"altron/config"
	"altron/core/dto"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/pkg/request"
	"altron/pkg/sftp"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.PcapWorkspaceUseCase = (*PcapWorkspaceUseCase)(nil)

type PcapWorkspaceUseCase struct {
	cfg               *config.AppConfig
	sftpClient        *sftp.Client
	pcapWorkspaceRepo interfaces.PcapWorkspaceRepository
}

func NewPcapWorkspaceUseCase(
	cfg *config.AppConfig,
	sftpClient *sftp.Client,
	pcapWorkspaceRepo interfaces.PcapWorkspaceRepository,
) *PcapWorkspaceUseCase {
	return &PcapWorkspaceUseCase{
		cfg:               cfg,
		sftpClient:        sftpClient,
		pcapWorkspaceRepo: pcapWorkspaceRepo,
	}
}

func (u *PcapWorkspaceUseCase) CreatePcapWorkspace(ctx context.Context, userID uuid.UUID, file multipart.File, fileHeader *multipart.FileHeader) (*spec.CreatePcapWorkspaceResponse, error) {
	if !strings.Contains(fileHeader.Filename, ".pcap") {
		return nil, models.ErrorOnlyPcapSupported
	}
	if strings.Contains(fileHeader.Filename, "../") {
		return nil, models.ErrorInvalidFilename
	}

	filepath := fmt.Sprintf("/files/%s", fileHeader.Filename)
	pcap, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filepath)

	if _, err := io.Copy(pcap, file); err != nil {
		return nil, err
	}
	pcap.Close()

	if err := u.sftpClient.Upload(filepath); err != nil {
		return nil, err
	}

	pcapWorkspaceEnt, err := u.pcapWorkspaceRepo.CreatePcapWorkspace(ctx, userID, &models.PcapWorkspace{
		FileName: fileHeader.Filename,
		Status:   common.LISTENING,
	})
	if err != nil {
		return nil, err
	}
	if err := request.Put(
		fmt.Sprintf("http://altron.connection.loc:%d/workspaces/pcaps/%s", u.cfg.AltronConnectionPort, pcapWorkspaceEnt.ID),
		dto.StartPcapListeningRequest{
			FileName: fileHeader.Filename,
		},
	); err != nil {
		return nil, err
	}
	statusCode, err := request.PostWithEmptyResponse(
		fmt.Sprintf("http://altron.session.loc:%d/pcaps", u.cfg.AltronSessionPort),
		dto.UploadPcapRequest{
			FileName: fileHeader.Filename,
		},
	)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, models.ErrorImportPcap
	}
	return &spec.CreatePcapWorkspaceResponse{
		Id:       pcapWorkspaceEnt.ID.String(),
		FileName: pcapWorkspaceEnt.FileName,
		Status:   spec.CreatePcapWorkspaceResponseStatus(pcapWorkspaceEnt.Status),
	}, nil
}

func (u *PcapWorkspaceUseCase) CreateLargePcapWorkspace(ctx context.Context, userID uuid.UUID, filename string) error {
	pcapWorkspaceEnt, err := u.pcapWorkspaceRepo.CreatePcapWorkspace(ctx, userID, &models.PcapWorkspace{
		FileName: filename,
		Status:   common.LISTENING,
	})
	if err != nil {
		return err
	}
	if err := request.Put(
		fmt.Sprintf("http://altron.connection.loc:%d/workspaces/pcaps/%s", u.cfg.AltronConnectionPort, pcapWorkspaceEnt.ID),
		dto.StartPcapListeningRequest{
			FileName: filename,
		},
	); err != nil {
		return err
	}
	statusCode, err := request.PostWithEmptyResponse(
		fmt.Sprintf("http://altron.session.loc:%d/pcaps",u.cfg.AltronSessionPort),
		dto.UploadPcapRequest{
			FileName: filename,
		},
	)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return models.ErrorImportPcap
	}
	return nil
}

func (u *PcapWorkspaceUseCase) UpdatePcapWorkspace(ctx context.Context, userID uuid.UUID, pcapWorkspaceID uuid.UUID, request *spec.UpdatePcapWorkspaceRequest) error {
	return u.pcapWorkspaceRepo.UpdatePcapWorkspace(ctx, userID, pcapWorkspaceID, &models.PcapWorkspace{
		Status: common.WorkspaceStatus(request.Status),
	})
}

func (u *PcapWorkspaceUseCase) GetPcapWorkspace(ctx context.Context, userID, pcapWorkspaceID uuid.UUID) (*spec.GetPcapWorkspaceResponse, error) {
	pcapWorkspace, err := u.pcapWorkspaceRepo.GetPcapWorkspaceByID(ctx, userID, pcapWorkspaceID)
	if err != nil {
		return nil, err
	}
	return &spec.GetPcapWorkspaceResponse{
		Id:       pcapWorkspace.ID.String(),
		FileName: pcapWorkspace.FileName,
		Status:   spec.GetPcapWorkspaceResponseStatus(pcapWorkspace.Status),
	}, nil
}

func (u *PcapWorkspaceUseCase) DeletePcapWorkspace(ctx context.Context, userID, pcapWorkspaceID uuid.UUID) error {
	return u.pcapWorkspaceRepo.DeletePcapWorkspace(ctx, userID, pcapWorkspaceID)
}
