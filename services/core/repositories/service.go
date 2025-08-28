package repositories

import (
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/plugin"
	"altron/core/repositories/ent/service"
	"altron/core/repositories/ent/user"
	"altron/core/repositories/ent/workspace"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.ServiceRepository = (*ServiceRepository)(nil)

type ServiceRepository struct {
	Client *ent.Client
}

func NewServiceRepository(client *ent.Client) *ServiceRepository {
	return &ServiceRepository{
		Client: client,
	}
}

func (r *ServiceRepository) CreateService(ctx context.Context, userID uuid.UUID, service *models.Service) (*ent.Service, error) {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	pluginsEnt := make([]*ent.Plugin, 0, len(service.Plugins))
	for _, p := range service.Plugins {
		pluginEnt, err := tx.Plugin.
			Query().
			Where(plugin.NameEQ(p.Name)).
			Only(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
		pluginsEnt = append(pluginsEnt, pluginEnt)
	}

	serviceEnt, err := tx.Service.
		Create().
		SetUserID(userID).
		SetName(service.Name).
		SetLink(service.Link).
		SetPort(uint32(service.Port)).
		SetNillableContainerID(service.ContainerID).
		AddPlugins(pluginsEnt...).
		Save(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		if strings.Contains(err.Error(), "duplicate") {
			return nil, models.ErrorServiceAlreadyExists
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return serviceEnt, nil
}

func (r *ServiceRepository) GetAllServices(ctx context.Context, userID uuid.UUID) ([]*ent.Service, error) {
	servicesEnt, err := r.Client.Service.
		Query().
		WithPlugins().
		Where(service.HasUserWith(user.IDEQ(userID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return servicesEnt, nil
}

func (r *ServiceRepository) GetServiceByID(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID) (*ent.Service, error) {
	serviceEnt, err := r.Client.Service.
		Query().
		WithPlugins().
		Where(service.And(
			service.HasUserWith(user.IDEQ(userID)),
			service.IDEQ(serviceID),
		)).
		Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, models.ErrorServiceNotFound
		}
		return nil, err
	}
	return serviceEnt, nil
}

func (r *ServiceRepository) GetAllPorts(ctx context.Context, userID uuid.UUID) ([]*ent.Service, error) {
	return r.Client.Service.
		Query().
		Select("port", "container_id").
		Where(service.HasUserWith(user.IDEQ(userID))).
		All(ctx)
}

func (r *ServiceRepository) UpdateService(ctx context.Context, userID uuid.UUID, serviceModel *models.Service) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}

	existingPlugins, err := tx.Service.
		Query().
		Where(service.And(
			service.HasUserWith(user.IDEQ(userID)),
			service.IDEQ(serviceModel.ID),
		)).
		QueryPlugins().
		All(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	pluginNames := make([]string, 0)
	for _, plugin := range serviceModel.Plugins {
		pluginNames = append(pluginNames, plugin.Name)
	}
	pluginsEnt, err := tx.Plugin.
		Query().
		Where(plugin.NameIn(pluginNames...)).
		All(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	qb := tx.Service.Update().
		SetNillableContainerID(serviceModel.ContainerID)

	if len(serviceModel.Name) != 0 {
		qb = qb.SetName(serviceModel.Name)
	}
	if len(serviceModel.Link) != 0 {
		qb = qb.SetLink(serviceModel.Link)
	}
	if len(existingPlugins) != 0 {
		qb = qb.RemovePlugins(existingPlugins...)
	}
	if len(pluginsEnt) != 0 {
		qb = qb.AddPlugins(pluginsEnt...)
	}
	err = qb.Where(service.And(
		service.HasUserWith(user.IDEQ(userID)),
		service.IDEQ(serviceModel.ID),
	)).Exec(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r *ServiceRepository) DeleteService(ctx context.Context, userID uuid.UUID, serviceID uuid.UUID) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Service.
		Delete().
		Where(service.And(
			service.HasUserWith(user.IDEQ(userID)),
			service.IDEQ(serviceID),
		)).
		Exec(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	_, err = tx.Workspace.
		Delete().
		Where(
			workspace.HasServiceWith(service.And(
				service.IDEQ(serviceID),
				service.HasUserWith(user.IDEQ(userID)),
			)),
		).
		Exec(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r *ServiceRepository) GetServicePluginsByPort(ctx context.Context, userID uuid.UUID, servicePort uint16) ([]*ent.Plugin, error) {
	plugins, err := r.Client.Plugin.
		Query().
		Where(plugin.HasServicesWith(
			service.And(
				service.HasUserWith(user.IDEQ(userID)),
				service.PortEQ(uint32(servicePort)),
			),
		)).
		All(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return make([]*ent.Plugin, 0), nil
		}
		return nil, err
	}
	return plugins, nil
}

func (r *ServiceRepository) GetServiceByPort(ctx context.Context, userID uuid.UUID, servicePort uint16) (*ent.Service, error) {
	service, err := r.Client.Service.
		Query().
		WithPlugins().
		WithWorkspaces().
		Where(service.And(
			service.HasUserWith(user.IDEQ(userID)),
			service.PortEQ(uint32(servicePort)),
		)).
		Only(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, models.ErrorServiceNotFound
		}
		return nil, err
	}
	return service, nil
}
