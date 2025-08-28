package filter

import (
	"altron/common/models"
	"altron/core/generated/spec"

	"github.com/google/uuid"
)

func PrepareUpdateFilterRequest(request *spec.UpdateFilterRequest, filterID uuid.UUID) *models.Filter {
	var serviceUuid *uuid.UUID
	if request.ServiceId != nil {
		id := uuid.MustParse(string(*request.ServiceId))
		serviceUuid = &id
	}

	return &models.Filter{
		ID:           filterID,
		Regex:        request.Regex,
		TTL:          request.Ttl,
		ServiceID:    serviceUuid,
		TotalPackets: request.TotalPackets,
		InRequest:    request.InRequest,
		InResponse:   request.InResponse,
		Color:        request.Color,
		IsBlocking:   request.IsBlocking,
	}
}

func PrepareToModels(filtersSpec []spec.Filter) []*models.Filter {
	filters := make([]*models.Filter, 0, len(filtersSpec))
	for _, filterSpec := range filtersSpec {
		filters = append(filters, &models.Filter{
			ID:           uuid.MustParse(filterSpec.Id),
			Name:         filterSpec.Name,
			Regex:        filterSpec.Regex,
			TTL:          filterSpec.Ttl,
			TotalPackets: filterSpec.TotalPackets,
			Color:        filterSpec.Color,
			InRequest:    filterSpec.InRequest,
			InResponse:   filterSpec.InResponse,
			IsBlocking:   filterSpec.IsBlocking,
		})
	}
	return filters
}
