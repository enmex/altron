package filter

import (
	"altron/core/generated/spec"
	"altron/core/repositories/ent"
)

func PresentFilters(filtersEnt []*ent.Filter) []spec.Filter {
	filters := make([]spec.Filter, 0, len(filtersEnt))
	for _, filterEnt := range filtersEnt {
		var regex *string
		var ttl *uint8
		var totalPackets *int
		if len(filterEnt.Regex) > 0 {
			regex = &filterEnt.Regex
		}
		if filterEnt.TTL != 0 {
			ttl = &filterEnt.TTL
		}
		if filterEnt.TotalPackets != 0 {
			totalPackets = &filterEnt.TotalPackets
		}
		var serviceID *string
		serviceEnt, err := filterEnt.Edges.ServiceOrErr()
		if err == nil {
			uuidStr := serviceEnt.ID.String()
			serviceID = &uuidStr
		}
		filters = append(filters, spec.Filter{
			Id:           filterEnt.ID.String(),
			Name:         filterEnt.Name,
			ServiceId:    serviceID,
			Regex:        regex,
			Ttl:          ttl,
			TotalPackets: totalPackets,
			Color:        filterEnt.Color,
			InRequest:    filterEnt.InRequest,
			InResponse:   filterEnt.InResponse,
			IsBlocking:   filterEnt.IsBlocking,
		})
	}
	return filters
}
