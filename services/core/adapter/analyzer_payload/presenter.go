package analyzerpayload

import (
	"altron/core/generated/spec"
	"altron/core/repositories/ent"
)

func PresentAnalyzerPayloads(analyzerPayload []*ent.AnalyzerPayload) []spec.Characteristic {
	resultPayload := make([]spec.Characteristic, 0)

	for _, characteristicEnt := range analyzerPayload {
		resultPayload = append(resultPayload, spec.Characteristic{
			Value:  characteristicEnt.Value,
			Number: characteristicEnt.Number,
		})
	}

	return resultPayload
}
