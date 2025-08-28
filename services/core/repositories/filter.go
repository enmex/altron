package repositories

import (
	common "altron/common/models"
	"altron/core/interfaces"
	"altron/core/models"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/filter"
	"altron/core/repositories/ent/service"
	"altron/core/repositories/ent/sessionfilter"
	"altron/core/repositories/ent/user"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ interfaces.FilterRepository = (*FilterRepository)(nil)

type FilterRepository struct {
	Client *ent.Client
}

func NewFilterRepository(client *ent.Client) *FilterRepository {
	return &FilterRepository{
		Client: client,
	}
}

func (r *FilterRepository) CreateFilter(ctx context.Context, userID uuid.UUID, filter *common.Filter) (*ent.Filter, error) {
	filterEnt, err := r.Client.Filter.
		Create().
		SetUserID(userID).
		SetInRequest(filter.InRequest).
		SetInResponse(filter.InResponse).
		SetNillableRegex(filter.Regex).
		SetNillableTTL(filter.TTL).
		SetNillableTotalPackets(filter.TotalPackets).
		SetNillableServiceID(filter.ServiceID).
		SetName(filter.Name).
		SetIsBlocking(filter.IsBlocking).
		SetColor(filter.Color).
		Save(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, models.ErrorDuplicateFilter
		}
		return nil, err
	}

	return filterEnt, nil
}

func (r *FilterRepository) GetAllFilters(ctx context.Context, userID uuid.UUID, servicePort uint16) ([]*ent.Filter, error) {
	return r.Client.Filter.
		Query().
		WithService().
		Where(filter.And(
			filter.Or(
				filter.HasServiceWith(service.PortEQ(uint32(servicePort))),
				filter.Not(filter.HasService()),
			),
			filter.HasUserWith(user.IDEQ(userID)),
		)).
		All(ctx)
}

func (r *FilterRepository) DeleteFilter(ctx context.Context, userID uuid.UUID, filterID uuid.UUID) (*ent.Filter, error) {
	filterEnt, err := r.Client.Filter.
		Query().
		Where(filter.And(
			filter.IDEQ(filterID), 
			filter.HasUserWith(user.IDEQ(userID)),
		)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.Client.Filter.
		Delete().
		Where(filter.IDEQ(filterID)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return filterEnt, nil
}

func (r *FilterRepository) UpdateFilter(ctx context.Context, userID uuid.UUID, filterModel *common.Filter) error {
	return r.Client.Filter.
		Update().
		SetNillableTTL(filterModel.TTL).
		SetNillableRegex(filterModel.Regex).
		SetNillableTotalPackets(filterModel.TotalPackets).
		SetNillableServiceID(filterModel.ServiceID).
		SetInRequest(filterModel.InRequest).
		SetInResponse(filterModel.InResponse).
		SetColor(filterModel.Color).
		SetIsBlocking(filterModel.IsBlocking).
		Where(filter.And(
			filter.IDEQ(filterModel.ID),
			filter.HasUserWith(user.IDEQ(userID)),
		)).
		Exec(ctx)
}

func (r *FilterRepository) CreateSessionFilters(ctx context.Context, sessionFilters []*common.SessionFilter) error {
	tx, err := r.Client.Tx(ctx)
	if err != nil {
		return err
	}

	for _, sessionFilter := range sessionFilters {
		sessionFilterEnt, err := tx.SessionFilter.
			Create().
			SetFilterID(sessionFilter.ID).
			SetSessionID(sessionFilter.SessionID).
			SetMatchesCount(sessionFilter.MatchesCount).
			Save(ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}

		for _, matchedPacketIdx := range sessionFilter.MatchedPackets {
			err := tx.MatchedPacket.
				Create().
				SetSessionFilter(sessionFilterEnt).
				SetPacketIdx(matchedPacketIdx).Exec(ctx)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *FilterRepository) DeleteSessionFilter(ctx context.Context, userID, filterID uuid.UUID) error {
	_, err := r.Client.SessionFilter.
		Delete().
		Where(sessionfilter.And(
			sessionfilter.HasFilterWith(filter.HasUserWith(user.IDEQ(userID))),
			sessionfilter.FilterIDEQ(filterID),
		)).Exec(ctx)
	return err
}