package handlers

import (
	analyzer "altron/common/analyzers"
	commonDto "altron/common/dto"
	"altron/common/interfaces"
	common "altron/common/models"
	"altron/config"
	"altron/pkg/redis"
	"altron/pkg/request"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type AnalyzerHandler struct {
	cfg          *config.AppConfig
	analyzers    map[string]interfaces.Analyzer
	checkerMasks map[uint16]*common.CheckerMask
	redisClient  *redis.Client[[]common.Characteristic]
}

func NewAnalyzerHandler(
	cfg *config.AppConfig,
	redisClient *redis.Client[[]common.Characteristic]) *AnalyzerHandler {
	analyzers := make(map[string]interfaces.Analyzer)
	analyzers["ttl"] = analyzer.NewTtlAnalyzer()
	analyzers["requests"] = analyzer.NewTotalRequestsAnalyzer()
	analyzers["ua"] = analyzer.NewUserAgentAnalyzer()
	analyzers["timestamps"] = analyzer.NewTimestampsAnalyzer()
	return &AnalyzerHandler{
		cfg:          cfg,
		analyzers:    analyzers,
		checkerMasks: make(map[uint16]*common.CheckerMask),
		redisClient:  redisClient,
	}
}

func (h *AnalyzerHandler) ServeSession(ctx context.Context, session *common.Session) (map[string]common.Characteristic, error) {
	analyzerMatches := make(map[string]common.Characteristic, 0)

	checkerMask, err := h.getServiceCheckerMask(ctx, session.ServerPort)
	if err != nil {
		return nil, err
	}
	isSafe := true

	for componentName, analyzer := range h.analyzers {
		value := analyzer.GetCharacteristicValue(ctx, session)
		if value != nil {
			if checkerMask.Present() {
				isSafe = isSafe && checkerMask.ContainsCharacteristic(componentName, *value)
			}
			key := fmt.Sprintf("%d_%s", session.ServerPort, componentName)
			var characteristics []common.Characteristic
			componentPayload, err := h.redisClient.Get(ctx, key)
			if err != nil {
				return nil, err
			}
			characteristics = *componentPayload

			found := false
			for _, characteristic := range characteristics {
				if strings.EqualFold(characteristic.Value, *value) {
					found = true
					break
				}
			}
			if !found {
				characteristics = append(characteristics, common.Characteristic{
					IsSafe: checkerMask.ContainsCharacteristic(componentName, *value),
					Value:  *value,
					Number: 0,
				})
			}

			if len(characteristics) != 0 {
				for idx := range characteristics {
					if strings.EqualFold(characteristics[idx].Value, *value) {
						characteristics[idx].Number++
						analyzerMatches[componentName] = characteristics[idx]
						break
					}
				}

				if err := h.redisClient.Set(ctx, key, characteristics); err != nil {
					return nil, err
				}
			}
		} else if checkerMask.ContainsComponent(componentName) {
			isSafe = false
		}
	}
	if checkerMask.Present() {
		session.IsSafe = &isSafe
	}
	return analyzerMatches, nil
}

func (h *AnalyzerHandler) PutWorkspaceAnalyzerMatches(ctx context.Context, workspaceID uuid.UUID, isChecker bool, servicePort uint16, analyzerMatches map[string]common.Characteristic) error {
	var workspaceId string
	if isChecker {
		workspaceId = "checker"
	} else {
		workspaceId = workspaceID.String()
	}

	for componentName, matchedCharacteristic := range analyzerMatches {
		key := fmt.Sprintf("%d_%s_%s", servicePort, workspaceId, componentName)

		var characteristics []common.Characteristic
		componentPayload, err := h.redisClient.Get(ctx, key)
		if err != nil {
			return err
		}
		characteristics = *componentPayload

		found := false
		for _, characteristic := range characteristics {
			if strings.EqualFold(characteristic.Value, matchedCharacteristic.Value) {
				found = true
				break
			}
		}
		if !found {
			newCharacteristic := common.Characteristic{
				Value:  matchedCharacteristic.Value,
				Number: 0,
			}
			characteristics = append(characteristics, newCharacteristic)
		}

		if len(characteristics) != 0 {
			for idx := range characteristics {
				if strings.EqualFold(characteristics[idx].Value, matchedCharacteristic.Value) {
					characteristics[idx].Number++
					break
				}
			}
			if err := h.redisClient.Set(ctx, key, characteristics); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *AnalyzerHandler) GetAnalyzerPayload(ctx context.Context, serverPort uint16, workspaceID *string) (*common.AnalyzerPayload, error) {
	analyzerPayload := make(map[string][]common.Characteristic, 0)

	for componentName := range h.analyzers {
		var key string
		if workspaceID != nil {
			key = fmt.Sprintf("%d_%s_%s", serverPort, *workspaceID, componentName)
		} else {
			key = fmt.Sprintf("%d_%s", serverPort, componentName)
		}

		var characteristics []common.Characteristic
		data, err := h.redisClient.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		characteristics = *data
		analyzerPayload[componentName] = characteristics
	}

	return &common.AnalyzerPayload{
		Payload: analyzerPayload,
	}, nil
}

func (h *AnalyzerHandler) MatchesSelectedCharacteristic(ctx context.Context, session *common.Session, selectedCharacteristics map[string][]commonDto.SelectedCharacteristicDto) bool {
	isMatch := true

	for componentName, selectedCharacteristics := range selectedCharacteristics {
		value := h.analyzers[componentName].GetCharacteristicValue(ctx, session)

		matches := false
		for _, selectedCharacteristic := range selectedCharacteristics {
			if strings.EqualFold(*value, selectedCharacteristic.Value) {
				if selectedCharacteristic.Blocked {
					return false
				}
				matches = true
				break
			}
		}
		if len(selectedCharacteristics) != 0 {
			isMatch = isMatch && matches
		}
	}

	return isMatch
}

func (h *AnalyzerHandler) ClearAnalyzerPayload(ctx context.Context, serverPort uint16, workspaceID *uuid.UUID) error {
	for componentName := range h.analyzers {
		var key string
		if workspaceID != nil {
			key = fmt.Sprintf("%d_%s_%s", serverPort, *workspaceID, componentName)
		} else {
			key = fmt.Sprintf("%d_%s", serverPort, componentName)
		}

		if err := h.redisClient.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func (h *AnalyzerHandler) getServiceCheckerMask(ctx context.Context, servicePort uint16) (*common.CheckerMask, error) {
	checkerMask, ok := h.checkerMasks[servicePort]

	if !ok {
		checkerMaskResponse, err := request.Get[commonDto.GetServiceCheckerMaskResponse](
			fmt.Sprintf("http://%s:%d/api/session-analyzer/services/%d/checker", h.cfg.AltronHost, h.cfg.AltronPort, servicePort),
		)
		if err != nil {
			return nil, err
		}
		checkerMask = &common.CheckerMask{
			AnalyzerPayload: checkerMaskResponse.Analyzer,
		}
		h.checkerMasks[servicePort] = checkerMask
	}

	return checkerMask, nil
}
