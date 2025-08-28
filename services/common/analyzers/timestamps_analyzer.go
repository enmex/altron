package analyzer

import (
	"altron/common/interfaces"
	common "altron/common/models"
	"context"
)

var _ interfaces.Analyzer = (*TimestampsAnalyzer)(nil)

type TimestampsAnalyzer struct{}

func NewTimestampsAnalyzer() *TimestampsAnalyzer {
	return &TimestampsAnalyzer{}
}

func (a *TimestampsAnalyzer) GetCharacteristicValue(ctx context.Context, session *common.Session) *string {
	avgTimestamps := session.AverageResponseTime

	var resStr string
	if avgTimestamps < 5 {
		resStr = "0-5ms"
	} else if avgTimestamps < 15 {
		resStr = "5-15ms"
	} else if avgTimestamps < 60 {
		resStr = "15-60ms"
	} else if avgTimestamps < 120 {
		resStr = "60-120ms"
	} else if avgTimestamps < 1000 {
		resStr = "120-1000ms"
	} else {
		resStr = ">1000ms"
	}
	return &resStr
}
