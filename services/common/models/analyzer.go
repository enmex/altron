package models

import (
	"altron/utils"
	"strings"
)

type Characteristic struct {
	Value  string `json:"value"`
	Number int    `json:"number"`
	IsSafe bool   `json:"isSafe"`
}

type AnalyzerPayload struct {
	Payload map[string][]Characteristic `json:"payload"`
}

type CheckerMask struct {
	AnalyzerPayload map[string][]Characteristic
}

func (cm *CheckerMask) ContainsCharacteristic(componentName string, chValue string) bool {
	chs, ok := cm.AnalyzerPayload[componentName]
	if !ok {
		return false
	}
	return utils.ContainsFunc[Characteristic](chs, func(c Characteristic) bool {
		return strings.EqualFold(c.Value, chValue)
	})
}

func (cm *CheckerMask) ContainsComponent(componentName string) bool {
	_, ok := cm.AnalyzerPayload[componentName]
	return ok
}

func (cm *CheckerMask) Present() bool {
	for _, chs := range cm.AnalyzerPayload {
		if len(chs) != 0 {
			return true
		}
	}
	return false
}
