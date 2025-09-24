package usecases

import (
	common "altron/common/models"
	"altron/config"
	"altron/core/adapter/plugin"
	"altron/core/adapter/service"
	"altron/core/dto"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/pkg/redis"
	req "altron/pkg/request"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.ServiceUseCase = (*ServiceUseCase)(nil)

type ServiceUseCase struct {
	cfg         *config.Config
	redisClient *redis.Client[[]spec.Characteristic]
	serviceRepo interfaces.ServiceRepository
}

func NewServiceUseCase(
	cfg *config.Config,
	redisClient *redis.Client[[]spec.Characteristic],
	repo interfaces.ServiceRepository,
) *ServiceUseCase {
	return &ServiceUseCase{
		cfg:         cfg,
		redisClient: redisClient,
		serviceRepo: repo,
	}
}

func (u *ServiceUseCase) CreateService(ctx context.Context, userID uuid.UUID, request *spec.CreateServiceRequest) (*spec.CreateServiceResponse, error) {
	plugins := make([]*common.Plugin, 0, len(request.Plugins))
	for _, pluginName := range request.Plugins {
		plugins = append(plugins, &common.Plugin{
			Name: pluginName,
		})
	}
	service := &models.Service{
		Name:        request.Name,
		Link:        request.Link,
		Port:        request.Port,
		ContainerID: request.ContainerID,
		Plugins:     plugins,
	}

	serviceEnt, err := u.serviceRepo.CreateService(ctx, userID, service)
	if err != nil {
		return nil, err
	}

	statusCode, err := req.PostWithEmptyResponse(
		fmt.Sprintf("http://%s:%d/ports", u.cfg.App.AltronHost, u.cfg.App.AltronSessionPort),
		dto.CreatePortsRequest{
			Ports: []dto.PortRequest{{
				Port:        service.Port,
				ContainerID: request.ContainerID,
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	if statusCode != 201 {
		return nil, models.ErrorUnableToPutPortOnSessionTree
	}

	for _, componentName := range common.ComponentNames {
		if err := u.redisClient.Set(ctx, fmt.Sprintf("%d_%s", service.Port, componentName), make([]spec.Characteristic, 0)); err != nil {
			return nil, err
		}
	}
	res := &spec.CreateServiceResponse{
		Id:      serviceEnt.ID.String(),
		Link:    serviceEnt.Link,
		Name:    serviceEnt.Name,
		Port:    uint16(serviceEnt.Port),
		Plugins: plugin.PreparePlugins(serviceEnt.Edges.Plugins),
	}

	data, err := json.Marshal(dto.AddPortRequest{
		Port:        request.Port,
		ContainerID: request.ContainerID,
	})
	if err != nil {
		return nil, err
	}
	if _, err := req.PostWithEmptyResponse(
		fmt.Sprintf("http://%s:%d/events", u.cfg.App.AltronHost, u.cfg.App.AltronConnectionPort),
		dto.CreateEventRequest{
			Type:            models.CreateService,
			Data:            data,
			WaitForResponse: true,
		},
	); err != nil && !strings.Contains(err.Error(), "agent") {
		return nil, err
	}

	return res, nil
}

func (u *ServiceUseCase) GetService(ctx context.Context, userID uuid.UUID, servicePort uint16) (*spec.GetServiceResponse, error) {
	service, err := u.serviceRepo.GetServiceByPort(ctx, userID, servicePort)
	if err != nil {
		return nil, err
	}
	plugins := make([]string, 0, len(service.Edges.Plugins))
	for _, pluginEnt := range service.Edges.Plugins {
		plugins = append(plugins, pluginEnt.Name)
	}
	workspaces := make([]spec.ServiceWorkspaceResponse, 0, len(service.Edges.Workspaces))
	for _, workspaceEnt := range service.Edges.Workspaces {
		workspaces = append(workspaces, spec.ServiceWorkspaceResponse{
			Id:     workspaceEnt.ID.String(),
			Name:   workspaceEnt.Name,
			Status: spec.ServiceWorkspaceResponseStatus(workspaceEnt.Status),
		})
	}
	return &spec.GetServiceResponse{
		Id:         service.ID.String(),
		Link:       service.Link,
		Name:       service.Name,
		Plugins:    plugins,
		Port:       uint16(service.Port),
		Workspaces: workspaces,
	}, nil
}

func (u *ServiceUseCase) GetAllPorts(ctx context.Context, userID uuid.UUID) (*dto.GetPortsResponse, error) {
	services, err := u.serviceRepo.GetAllPorts(ctx, userID)
	if err != nil {
		return nil, err
	}

	ports := make([]dto.PortRequest, 0)
	for _, service := range services {
		ports = append(ports, dto.PortRequest{
			Port:        uint16(service.Port),
			ContainerID: service.ContainerID,
		})
	}

	return &dto.GetPortsResponse{
		Ports: ports,
	}, nil
}

func (u *ServiceUseCase) UpdateService(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID, request *spec.UpdateServiceRequest) error {
	serviceEnt, err := u.serviceRepo.GetServiceByID(ctx, userID, serviceID)
	if err != nil {
		return err
	}
	if request.Name == nil && request.Link == nil {
		return models.ErrorAllFieldsEmpty
	}
	if err := u.serviceRepo.UpdateService(ctx, userID, service.PrepareUpdateServiceRequest(request, serviceID)); err != nil {
		return err
	}
	if request.ContainerID == nil {
		return nil
	}
	_, err = req.PatchWithEmptyResponse(fmt.Sprintf("http://%s:%d/logs", u.cfg.App.AltronHost, u.cfg.App.AltronSessionPort), dto.UpdateContainerRequest{
		OldContainer: serviceEnt.ContainerID,
		NewContainer: *request.ContainerID,
	})
	if err != nil {
		return err
	}
	data, err := json.Marshal(dto.PortRequest{
		Port:        uint16(serviceEnt.Port),
		ContainerID: request.ContainerID,
	})
	if err != nil {
		return err
	}
	res, err := req.Post[dto.CreateEventResponse](
		fmt.Sprintf("http://altron.connection.loc:%d/events", u.cfg.App.AltronConnectionPort),
		dto.CreateEventRequest{
			Type:            models.UpdateService,
			Data:            data,
			WaitForResponse: true,
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "agent") {
			return nil
		}
		return err
	}
	if res.Status != "ok" {
		return fmt.Errorf("%v", res.Data)
	}
	return nil
}

func (u *ServiceUseCase) DeleteService(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID) error {
	service, err := u.serviceRepo.GetServiceByID(ctx, userID, serviceID)
	if err != nil {
		return err
	}

	if err := u.serviceRepo.DeleteService(ctx, userID, serviceID); err != nil {
		return err
	}

	statusCode, err := req.DeleteWithRequest(fmt.Sprintf("http://%s:%d/ports", u.cfg.App.AltronHost, u.cfg.App.AltronSessionPort), dto.DeletePortRequest{
		PortRequest: dto.PortRequest{
			Port:        uint16(service.Port),
			ContainerID: service.ContainerID,
		},
	})
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return models.ErrorUnableToDeletePortOnSessionTree
	}

	for _, componentName := range common.ComponentNames {
		if err := u.redisClient.Delete(ctx, fmt.Sprintf("%d_%s", service.Port, componentName)); err != nil {
			return err
		}
	}

	data, err := json.Marshal(dto.DeletePortRequest{
		PortRequest: dto.PortRequest{
			Port:        uint16(service.Port),
			ContainerID: service.ContainerID,
		},
	})
	if err != nil {
		return err
	}
	res, err := req.Post[dto.CreateEventResponse](
		fmt.Sprintf("http://%s:%d/events", u.cfg.App.AltronHost, u.cfg.App.AltronConnectionPort),
		dto.CreateEventRequest{
			Type:            models.DeleteService,
			Data:            data,
			WaitForResponse: true,
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "agent") {
			return nil
		}
		return err
	}
	if res.Status != "ok" {
		return fmt.Errorf("%v", res.Data)
	}

	return nil
}

func (u *ServiceUseCase) ScanHost(ctx context.Context, scope string) (*spec.ScanHostServicesResponse, error) {
	var eventType models.EventType
	if scope == "containers" {
		eventType = models.ScanContainers
	} else {
		eventType = models.ScanHost
	}

	res, err := req.Post[dto.CreateEventResponse](
		fmt.Sprintf("http://altron.connection.loc:%d/events", u.cfg.App.AltronConnectionPort),
		dto.CreateEventRequest{
			Type:            eventType,
			WaitForResponse: true,
		},
	)
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, fmt.Errorf("%v", res.Data)
	}
	var scanResponse spec.ScanHostServicesResponse
	if err := json.Unmarshal(res.Data, &scanResponse); err != nil {
		return nil, err
	}

	return &scanResponse, nil
}

func (u *ServiceUseCase) ScanService(ctx context.Context, port uint16) (*spec.ScanHostServiceResponse, error) {
	data, err := json.Marshal(dto.ScanHostServiceRequest{
		Port: port,
	})
	if err != nil {
		return nil, err
	}
	res, err := req.Post[dto.CreateEventResponse](
		fmt.Sprintf("http://altron.connection.loc:%d/events", u.cfg.App.AltronConnectionPort),
		dto.CreateEventRequest{
			Type:            models.ScanPort,
			Data:            data,
			WaitForResponse: true,
		},
	)
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, fmt.Errorf("%v", res.Data)
	}
	var scanResponse spec.ScanHostServiceResponse
	if err := json.Unmarshal(res.Data, &scanResponse); err != nil {
		return nil, err
	}

	return &scanResponse, nil
}
