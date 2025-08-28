package interfaces

import (
	common "altron/common/models"
	"context"
)

type Analyzer interface {
	GetCharacteristicValue(ctx context.Context, session *common.Session) *string
}
