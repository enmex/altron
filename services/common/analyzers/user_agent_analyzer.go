package analyzer

import (
	"altron/common/interfaces"
	common "altron/common/models"
	"context"
)

var _ interfaces.Analyzer = (*UserAgentAnalyzer)(nil)

type UserAgentAnalyzer struct{}

func NewUserAgentAnalyzer() *UserAgentAnalyzer {
	return &UserAgentAnalyzer{}
}

func (a *UserAgentAnalyzer) GetCharacteristicValue(ctx context.Context, session *common.Session) *string {
	return session.ClientUserAgent
}
