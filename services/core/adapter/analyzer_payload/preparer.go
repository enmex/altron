package analyzerpayload

import (
	common "altron/common/models"
	"altron/core/generated/spec"
)

func PrepareCreateAnalyzerPayloadRequest(request *spec.CreateAnalyzerPayloadRequest) *common.AnalyzerPayload {
	payload := make(map[string][]common.Characteristic, 0)

	for componentName, characteristicsSpec := range request.Payload.AdditionalProperties {
		characteristics := make([]common.Characteristic, 0, len(characteristicsSpec))

		for _, characteristicSpec := range characteristicsSpec {
			characteristics = append(characteristics, common.Characteristic{
				Value:  characteristicSpec.Value,
				Number: characteristicSpec.Number,
			})
		}
		payload[componentName] = characteristics
	}

	return &common.AnalyzerPayload{
		Payload: payload,
	}
}
