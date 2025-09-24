package usecases

import (
	commonHandlers "altron/common/handlers"
	common "altron/common/models"
	"altron/config"
	"altron/connection/dto"
	"altron/connection/generated/spec"
	"altron/connection/handlers"
	"altron/connection/interfaces"
	"altron/connection/models"
	"altron/pkg/amqp"
	req "altron/pkg/request"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/procyon-projects/chrono"
	"github.com/sirupsen/logrus"
)

var _ interfaces.WorkspaceUseCase = (*WorkspaceUseCase)(nil)

type WorkspaceUseCase struct {
	log             *logrus.Logger
	cfg             *config.AppConfig
	sessionHandler  *handlers.SessionHandler
	analyzerHandler *commonHandlers.AnalyzerHandler
	scheduler       chrono.TaskScheduler
	client          *amqp.Client
	currentTasks    map[uuid.UUID]chrono.ScheduledTask
	closeChans      map[uuid.UUID]chan bool
}

func NewWorkspaceUseCase(
	log *logrus.Logger,
	cfg *config.AppConfig,
	client *amqp.Client,
	sessionHandler *handlers.SessionHandler,
	analyzerHandler *commonHandlers.AnalyzerHandler,
) *WorkspaceUseCase {
	return &WorkspaceUseCase{
		log:             log,
		cfg:             cfg,
		sessionHandler:  sessionHandler,
		analyzerHandler: analyzerHandler,
		scheduler:       chrono.NewDefaultTaskScheduler(),
		client:          client,
		currentTasks:    make(map[uuid.UUID]chrono.ScheduledTask),
		closeChans:      make(map[uuid.UUID]chan bool),
	}
}

func (u *WorkspaceUseCase) StartListeningSessions(ctx context.Context, workspaceID uuid.UUID, servicePort uint16, request *spec.StartSessionsListeningRequest) error {
	timeout, err := time.ParseDuration(request.Timeout)
	if err != nil {
		return err
	}

	var startTime time.Time
	if request.StartTime == nil {
		startTime = time.Now()
	} else {
		startTime = time.UnixMilli(*request.StartTime)
	}
	u.closeChans[workspaceID] = make(chan bool)

	task, err := u.scheduler.Schedule(func(ctx context.Context) {
		if err := u.listenToSessions(ctx, workspaceID, request.IsChecker, servicePort, timeout); err != nil {
			u.log.Errorln(err)
		}
	}, chrono.WithTime(startTime))
	if err != nil {
		return err
	}
	u.currentTasks[workspaceID] = task
	return nil
}

func (u *WorkspaceUseCase) listenToSessions(ctx context.Context, workspaceID uuid.UUID, isChecker bool, servicePort uint16, timeout time.Duration) error {
	u.log.Infof("workspace-%s has started listening to port %d", workspaceID, servicePort)
	defer u.log.Infof("workspace-%s has disconnected from port %d", workspaceID, servicePort)

	consumer, err := amqp.NewConsumer(u.client, uuid.New().String())
	if err != nil {
		return err
	}
	defer consumer.Close()

	if err := consumer.Bind(fmt.Sprint(servicePort)); err != nil {
		return err
	}
	_, err = req.PatchWithEmptyResponse(
		fmt.Sprintf("http://altron.core.loc:%d/api/workspaces/%s", u.cfg.AltronPort, workspaceID),
		dto.UpdateWorkspaceStatusRequest{
			Status: common.LISTENING,
		})
	if err != nil {
		u.log.Errorln(err)
	}
	closeChan := u.closeChans[workspaceID]
	pluginsResponse, err := req.Get[dto.GetServicePluginsResponse](
		fmt.Sprintf("http://altron.core.loc:%d/api/services/%d", u.cfg.AltronPort, servicePort),
	)
	if err != nil {
		return err
	}

	filtersResponse, err := req.Get[dto.GetAllFiltersResponse](
		fmt.Sprintf("http://altron.core.loc:%d/api/filters?servicePort=%d", u.cfg.AltronPort, servicePort),
	)
	if err != nil {
		return err
	}

	sessionChan, err := consumer.Messages()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeout)
	for {
		select {
		case <-ticker.C:
			_, err := req.PatchWithEmptyResponse(
				fmt.Sprintf("http://altron.core.loc:%d/api/workspaces/%s", u.cfg.AltronPort, workspaceID),
				dto.UpdateWorkspaceStatusRequest{
					Status: common.COMPLETED,
				})
			if err != nil {
				return err
			}
			return u.saveAnalyzerPayload(ctx, servicePort, workspaceID, isChecker)
		case message := <-sessionChan:
			if err := u.processMessage(ctx, message.Body, workspaceID, isChecker, filtersResponse.Filters, pluginsResponse.Plugins); err != nil {
				u.log.Errorln(err)
			}
		case <-closeChan:
			return u.saveAnalyzerPayload(ctx, servicePort, workspaceID, isChecker)
		}
	}
}

func (u *WorkspaceUseCase) processMessage(ctx context.Context, body []byte, workspaceID uuid.UUID, isChecker bool, filters []*common.Filter, plugins []string) error {
	var session common.Session
	if err := json.Unmarshal(body, &session); err != nil {
		return err
	}

	processedSession, err := u.sessionHandler.ServeSession(&session, filters, plugins)
	if err != nil {
		return err
	}
	if processedSession == nil {
		return nil
	}
	if err := u.analyzerHandler.PutWorkspaceAnalyzerMatches(ctx, workspaceID, isChecker, session.ServerPort, session.AnalyzerMatches); err != nil {
		return err
	}
	statusCode, err := req.PostWithEmptyResponse(
		fmt.Sprintf("http://altron.core.loc:%d/api/workspaces/%s/sessions",
			u.cfg.AltronPort,
			workspaceID,
		), dto.AddSessionsToWorkspaceRequest{
			Sessions: []common.Session{*processedSession},
		})
	if err != nil {
		u.log.Errorln(err)
	}
	if statusCode != 200 {
		u.log.Errorln(models.ErrorAddingSessionsToWorkspace)
	}
	return nil
}

func (u *WorkspaceUseCase) saveAnalyzerPayload(ctx context.Context, servicePort uint16, workspaceID uuid.UUID, isChecker bool) error {
	var workspaceId string
	if isChecker {
		workspaceId = "checker"
	} else {
		workspaceId = workspaceID.String()
	}
	analyzerPayload, err := u.analyzerHandler.GetAnalyzerPayload(ctx, servicePort, &workspaceId)
	if err != nil {
		return err
	}
	statusCode, err := req.PostWithEmptyResponse(fmt.Sprintf("http://altron.core.loc:%d/api/session-analyzer/workspaces/%s",
		u.cfg.AltronPort,
		workspaceID,
	), dto.CreateAnalyzerPayloadRequest{
		Payload: analyzerPayload.Payload,
	})
	if err != nil {
		return err
	}
	if statusCode != http.StatusCreated {
		return models.ErrorCreatingSessionAnalyzerData
	}
	return u.analyzerHandler.ClearAnalyzerPayload(ctx, servicePort, &workspaceID)
}

func (u *WorkspaceUseCase) StopListeningSessions(workspaceID uuid.UUID) {
	u.currentTasks[workspaceID].Cancel()
	u.closeChans[workspaceID] <- true
	delete(u.currentTasks, workspaceID)
	delete(u.closeChans, workspaceID)
}
