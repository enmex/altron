package analyzer

import (
	"altron/common/interfaces"
	common "altron/common/models"
	"context"
	"fmt"
)

var _ interfaces.Analyzer = (*TotalRequestsAnalyzer)(nil)

type TotalRequestsAnalyzer struct{}

func NewTotalRequestsAnalyzer() *TotalRequestsAnalyzer {
	return &TotalRequestsAnalyzer{}
}

func (a *TotalRequestsAnalyzer) GetCharacteristicValue(ctx context.Context, session *common.Session) *string {
	var requestsRange string
	requestsNumber := session.RequestsNumber

	if requestsNumber > 15 && requestsNumber <= 20 {
		requestsRange = "15-20"
	} else if requestsNumber > 20 && requestsNumber <= 30 {
		requestsRange = "20-30"
	} else if requestsNumber > 30 {
		requestsRange = ">30"
	} else {
		requestsRange = fmt.Sprint(requestsNumber)
	}
	return &requestsRange
}
