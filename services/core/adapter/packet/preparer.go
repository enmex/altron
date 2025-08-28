package packet

import (
	"altron/common/models"
	"altron/core/generated/spec"
)

func PrepareConvertPacketToExploitRequest(request *spec.ConvertPacketToExploitRequest) *models.Packet {
	packet := request.Packet

	return &models.Packet{
		SentAt:    packet.SentAt,
		Payload:   packet.Payload,
		IsRequest: packet.IsRequest,
	}
}
