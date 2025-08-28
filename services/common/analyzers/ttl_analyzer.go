package analyzer

import (
	"altron/common/interfaces"
	common "altron/common/models"
	"context"
	"fmt"
)

var _ interfaces.Analyzer = (*TtlAnalyzer)(nil)

type TtlAnalyzer struct{}

func NewTtlAnalyzer() *TtlAnalyzer {
	return &TtlAnalyzer{}
}

func (a *TtlAnalyzer) GetCharacteristicValue(ctx context.Context, session *common.Session) *string {
	ttl := fmt.Sprint(session.TTL)
	return &ttl
}
