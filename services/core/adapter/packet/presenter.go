package packet

import (
	"altron/common/models"
	"altron/core/repositories/ent"
)

func PresentPacketsEnt(packetsEnt []*ent.Packet) []*models.Packet {
	packetsSpec := make([]*models.Packet, 0, len(packetsEnt))
	for _, packetEnt := range packetsEnt {
		packetsSpec = append(packetsSpec, &models.Packet{
			IsRequest: packetEnt.IsRequest,
			Payload:   packetEnt.Payload,
			SentAt:    packetEnt.SentAt,
		})
	}
	return packetsSpec
}
