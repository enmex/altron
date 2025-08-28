package usecases

import (
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"context"

	"github.com/google/uuid"
)

var _ interfaces.DashboardUseCase = (*DashboardUseCase)(nil)

type DashboardUseCase struct {
	serviceRepo       interfaces.ServiceRepository
	pcapWorkspaceRepo interfaces.PcapWorkspaceRepository
}

func NewDashboardUseCase(
	serviceRepo interfaces.ServiceRepository,
	pcapWorkspaceRepo interfaces.PcapWorkspaceRepository,
) *DashboardUseCase {
	return &DashboardUseCase{
		serviceRepo:       serviceRepo,
		pcapWorkspaceRepo: pcapWorkspaceRepo,
	}
}

func (u *DashboardUseCase) GetDashboard(ctx context.Context, userID uuid.UUID) (*spec.GetDashboardResponse, error) {
	services, err := u.serviceRepo.GetAllServices(ctx, userID)
	if err != nil {
		return nil, err
	}
	pcapWorkspaces, err := u.pcapWorkspaceRepo.GetAllPcapWorkspaces(ctx, userID)
	if err != nil {
		return nil, err
	}
	servicesSpec := make([]spec.DashboardServiceResponse, 0, len(services))
	for _, service := range services {
		servicesSpec = append(servicesSpec, spec.DashboardServiceResponse{
			Name:        service.Name,
			Port:        uint16(service.Port),
			ContainerID: service.ContainerID,
		})
	}
	pcapWorkspacesSpec := make([]spec.PcapWorkspace, 0, len(pcapWorkspaces))
	for _, pcapWorkspace := range pcapWorkspaces {
		pcapWorkspacesSpec = append(pcapWorkspacesSpec, spec.PcapWorkspace{
			Id:       pcapWorkspace.ID.String(),
			FileName: pcapWorkspace.FileName,
		})
	}
	return &spec.GetDashboardResponse{
		PcapWorkspaces: pcapWorkspacesSpec,
		Services:       servicesSpec,
	}, nil
}
