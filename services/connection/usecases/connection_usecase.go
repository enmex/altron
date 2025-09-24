package usecases

import (
	commonDto "altron/common/dto"
	common "altron/common/models"
	"altron/config"
	"altron/connection/dto"
	"altron/connection/generated/spec"
	"altron/connection/handlers"
	"altron/connection/interfaces"
	"altron/connection/models"
	"altron/pkg/amqp"
	"altron/pkg/request"
	"altron/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var _ interfaces.ConnectionUseCase = (*ConnectionUseCase)(nil)

type ConnectionUseCase struct {
	log                *logrus.Logger
	cfg                *config.AppConfig
	client             *amqp.Client
	sessionHandler     *handlers.SessionHandler
	eventsChan         chan *spec.CreateEventRequest
	eventsResponseChan chan *spec.CreateEventResponse
	mainChannelOpened  bool
	statsChannelOpened bool
	mut                *sync.Mutex
}

func NewConnectionUseCase(
	log *logrus.Logger,
	cfg *config.AppConfig,
	client *amqp.Client,
	sessionHandler *handlers.SessionHandler,
) *ConnectionUseCase {
	return &ConnectionUseCase{
		log:                log,
		cfg:                cfg,
		client:             client,
		sessionHandler:     sessionHandler,
		eventsChan:         make(chan *spec.CreateEventRequest),
		eventsResponseChan: make(chan *spec.CreateEventResponse),
		mainChannelOpened:  false,
		statsChannelOpened: false,
		mut:                &sync.Mutex{},
	}
}

func (u *ConnectionUseCase) ListenSessions(ctx context.Context, ws *websocket.Conn, servicePort uint16) error {
	consumer, err := amqp.NewConsumer(u.client, uuid.New().String())
	if err != nil {
		return err
	}
	defer consumer.Close()
	if err := consumer.Bind(fmt.Sprint(servicePort)); err != nil {
		return err
	}
	metricsConsumer, err := amqp.NewConsumer(u.client, uuid.New().String())
	if err != nil {
		return err
	}
	defer metricsConsumer.Close()
	if err := metricsConsumer.Bind(fmt.Sprintf("metrics-%d", servicePort)); err != nil {
		return err
	}

	closeChan := make(chan bool)
	errChan := make(chan error)
	isPause := false

	analyzerBlackList := make(map[string][]string, 0)
	analyzerPassList := make(map[string][]string, 0)

	pluginsResponse, err := request.Get[dto.GetServicePluginsResponse](
		fmt.Sprintf("http://altron.core.loc:%d/api/services/%d", u.cfg.AltronPort, servicePort),
	)
	if err != nil {
		return err
	}

	filtersResponse, err := request.Get[dto.GetAllFiltersResponse](
		fmt.Sprintf("http://altron.core.loc:%d/api/filters?servicePort=%d", u.cfg.AltronPort, servicePort),
	)
	if err != nil {
		return err
	}

	u.log.Infof("%s has connected to port %d", consumer.MemberID(), servicePort)

	defer func() {
		u.log.Infof("%s has disconnected from port %d", consumer.MemberID(), servicePort)
		consumer.Close()
	}()

	metricsChan, err := metricsConsumer.Messages()
	if err != nil {
		return err
	}

	if err := ws.WriteMessage(websocket.TextMessage, []byte("connect")); err != nil {
		return err
	}

	go func(ctx context.Context, ws *websocket.Conn) {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					closeChan <- true
					return
				}
				errChan <- err
				return
			}

			var clientRequest dto.ClientDto
			if err := json.Unmarshal(data, &clientRequest); err != nil {
				errChan <- err
				return
			}
			switch clientRequest.Action {
			case "refresh":
				u.log.Infof("%s refreshed session streaming in port %d", consumer.MemberID(), servicePort)
				filtersResponse, err = request.Get[dto.GetAllFiltersResponse](
					fmt.Sprintf("http://altron.core.loc:%d/api/filters?servicePort=%d", u.cfg.AltronPort, servicePort),
				)
				if err != nil {
					errChan <- err
					return
				}
				pluginsResponse, err = request.Get[dto.GetServicePluginsResponse](
					fmt.Sprintf("http://altron.core.loc:%d/api/services/%d", u.cfg.AltronPort, servicePort),
				)
				if err != nil {
					errChan <- err
					return
				}
			case "pause":
				u.log.Infof("%s paused session streaming in port %d", consumer.MemberID(), servicePort)
				isPause = !isPause
			case "analyzer-block":
				key, value := clientRequest.Payload.Get()
				u.log.Infof("%s blocked sessions with %s:%s in port %d", consumer.MemberID(), key, value, servicePort)
				characteristics, ok := analyzerBlackList[key]
				if !ok {
					analyzerBlackList[key] = make([]string, 0)
					analyzerBlackList[key] = append(analyzerBlackList[key], value)
				} else {
					updated := utils.ToggleElement[string](characteristics, value, func(i, j string) bool {
						return strings.EqualFold(i, j)
					})
					if len(updated) > 0 {
						analyzerBlackList[key] = updated
					} else {
						delete(analyzerBlackList, key)
					}
				}
			case "analyzer-pass":
				key, value := clientRequest.Payload.Get()
				u.log.Infof("%s filtered sessions with %s:%s in port %d", consumer.MemberID(), key, value, servicePort)
				characteristics, ok := analyzerPassList[key]
				if !ok {
					analyzerPassList[key] = make([]string, 0)
					analyzerPassList[key] = append(analyzerPassList[key], value)
				} else {
					updated := utils.ToggleElement[string](characteristics, value, func(i, j string) bool {
						return strings.EqualFold(i, j)
					})
					if len(updated) > 0 {
						analyzerPassList[key] = updated
					} else {
						delete(analyzerPassList, key)
					}
				}
			}
		}
	}(ctx, ws)

	sessionChan, err := consumer.Messages()
	if err != nil {
		return err
	}

	for {
		select {
		case message := <-sessionChan:
			if isPause {
				continue
			}
			var session common.Session
			if err := json.Unmarshal(message.Body, &session); err != nil {
				return err
			}

			processedSession, err := u.sessionHandler.ServeSession(&session, filtersResponse.Filters, pluginsResponse.Plugins)
			if err != nil {
				return err
			}
			if processedSession != nil {
				analyzerMatches := session.AnalyzerMatches

				match := true
				for componentName, characteristics := range analyzerPassList {
					characteristic, ok := analyzerMatches[componentName]
					if !ok {
						match = false
						break
					}
					if !utils.ContainsFunc[string](characteristics, func(s string) bool {
						return strings.EqualFold(s, characteristic.Value)
					}) {
						match = false
						break
					}
				}
				for componentName, characteristics := range analyzerBlackList {
					characteristic, ok := analyzerMatches[componentName]
					if ok && utils.ContainsFunc[string](characteristics, func(s string) bool {
						return strings.EqualFold(s, characteristic.Value)
					}) {
						match = false
						break
					}
				}

				if match {
					sessionJson, err := json.Marshal(session)
					if err != nil {
						return err
					}
					if err := ws.WriteMessage(websocket.TextMessage, sessionJson); err != nil {
						if !strings.Contains(err.Error(), "close sent") {
							return err
						}
						return nil
					}
				}
			}
		case metrics := <-metricsChan:
			if err := ws.WriteMessage(websocket.TextMessage, metrics.Body); err != nil {
				if !strings.Contains(err.Error(), "close sent") {
					return err
				}
				return nil
			}
		case <-closeChan:
			return nil
		case err := <-errChan:
			return err
		}
	}
}

func (u *ConnectionUseCase) ListenLogs(ctx context.Context, ws *websocket.Conn, containerID string) error {
	consumer, err := amqp.NewConsumer(u.client, uuid.New().String())
	if err != nil {
		return err
	}
	defer consumer.Close()
	if err := consumer.Bind(fmt.Sprintf("logs-%s", containerID)); err != nil {
		return err
	}

	u.log.Infof("%s has connected to listen logs of container %s", consumer.MemberID(), containerID)
	defer u.log.Infof("%s has stopped listening logs of container %s", consumer.MemberID(), containerID)

	if err := ws.WriteMessage(websocket.TextMessage, []byte("connect")); err != nil {
		return err
	}

	closeChan := make(chan bool)
	errChan := make(chan error)
	isPause := false

	go func(ctx context.Context, ws *websocket.Conn, containerID string) {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					closeChan <- true
					return
				}
				errChan <- err
				return
			}

			var clientRequest dto.ClientDto
			if err := json.Unmarshal(data, &clientRequest); err != nil {
				errChan <- err
				return
			}
			switch clientRequest.Action {
			case "pause":
				u.log.Infof("%s paused logs streaming in container %s", consumer.MemberID(), containerID)
				isPause = !isPause
			}
		}
	}(ctx, ws, containerID)

	logsDtoChan, err := consumer.Messages()
	if err != nil {
		return err
	}

	for {
		select {
		case message := <-logsDtoChan:
			if isPause {
				continue
			}
			if err := ws.WriteMessage(websocket.TextMessage, message.Body); err != nil {
				if !strings.Contains(err.Error(), "close sent") {
					return err
				}
				return nil
			}
		case <-closeChan:
			return nil
		case err := <-errChan:
			return err
		}
	}
}

func (u *ConnectionUseCase) ListenStats(ctx context.Context, ws *websocket.Conn) error {
	consumer, err := amqp.NewConsumer(u.client, uuid.New().String())
	if err != nil {
		return err
	}
	defer consumer.Close()
	if err := consumer.Bind("stats"); err != nil {
		return err
	}

	u.log.Infof("%s has connected to listen to altron`s health", consumer.MemberID())
	defer u.log.Infof("%s has stopped listen to altron`s health", consumer.MemberID())

	if err := ws.WriteMessage(websocket.TextMessage, []byte("connect")); err != nil {
		return err
	}

	closeChan := make(chan bool)
	errChan := make(chan error)

	go func(ctx context.Context, ws *websocket.Conn) {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					closeChan <- true
					return
				}
				errChan <- err
				return
			}
		}
	}(ctx, ws)

	statsChan, err := consumer.Messages()
	if err != nil {
		return err
	}

	for {
		select {
		case message := <-statsChan:
			if err := ws.WriteMessage(websocket.TextMessage, message.Body); err != nil {
				if !strings.Contains(err.Error(), "close sent") {
					return err
				}
				return nil
			}
		case <-closeChan:
			return nil
		case err := <-errChan:
			return err
		}
	}
}

func (u *ConnectionUseCase) OpenMainChannel(ctx context.Context, ws *websocket.Conn) error {
	u.log.Infoln("Main channel is opened")
	defer u.log.Infoln("Main channel is closed")

	producer, err := amqp.NewProducer(u.client)
	if err != nil {
		return err
	}
	defer producer.Close()

	if u.mainChannelOpened {
		data, err := json.Marshal(commonDto.ErrorResponse{
			Message: models.ErrorAgentAlreadyConnected.Error(),
		})
		if err != nil {
			return err
		}
		return ws.WriteJSON(spec.CreateEventRequest{
			Type: "error",
			Data: data,
		})
	}
	u.mut.Lock()
	u.mainChannelOpened = true
	u.mut.Unlock()
	defer func() {
		u.mut.Lock()
		u.mainChannelOpened = false
		u.mut.Unlock()
	}()

	//send services
	res, err := request.Get[dto.GetAllServicesResponse](
		fmt.Sprintf("http://altron.core.loc:%d/api/dashboard", u.cfg.AltronPort),
	)
	if err != nil {
		return err
	}
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	if err := ws.WriteJSON(spec.CreateEventRequest{
		Type: "get_dashboard",
		Data: data,
	}); err != nil {
		return err
	}

	//receive agent info
	var agentResponse spec.CreateEventResponse
	if err := ws.ReadJSON(&agentResponse); err != nil {
		return err
	}
	if agentResponse.Status == "error" {
		return fmt.Errorf("%s", string(agentResponse.Data))
	}
	var agentInfo dto.AgentInfoResponse
	if err := json.Unmarshal(agentResponse.Data, &agentInfo); err != nil {
		return err
	}
	if err := producer.CreateExchange(fmt.Sprintf("logs-%s", agentInfo.ContainerID)); err != nil {
		return err
	}
	defer producer.DeleteExchange(fmt.Sprintf("logs-%s", agentInfo.ContainerID))

	closeChan := make(chan bool)
	defer func() {
		closeChan <- true
	}()

	go func() {
		for {
			select {
			case req := <-u.eventsChan:
				bytes, err := json.Marshal(req)
				if err != nil {
					u.log.Errorln(err)
					return
				}
				if err := ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
					u.log.Errorln("error writing message:", err)
					return
				}
			case <-closeChan:
				return
			}
		}
	}()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("error reading message: %v", err)
		}
		if messageType == websocket.CloseMessage {
			return nil
		}
		u.log.Infoln("event response received")
		var response spec.CreateEventResponse
		if err := json.Unmarshal(message, &response); err != nil {
			return err
		}
		u.eventsResponseChan <- &response
	}
}

func (u *ConnectionUseCase) CreateEvent(ctx context.Context, request *spec.CreateEventRequest) (*spec.CreateEventResponse, error) {
	u.mut.Lock()
	if !u.mainChannelOpened {
		u.mut.Unlock()
		return nil, models.ErrorAgentNotConnected
	}
	u.mut.Unlock()
	u.eventsChan <- request
	if !request.WaitForResponse {
		return nil, nil
	}
	res, ok := <-u.eventsResponseChan
	u.log.Infoln(request.Type, *res)
	if !ok {
		return nil, models.ErrorSomethingWrong
	}
	return res, nil
}

func (u *ConnectionUseCase) GetAgentStatus(ctx context.Context) *spec.GetAgentStatusResponse {
	u.mut.Lock()
	defer u.mut.Unlock()

	var status string
	if u.mainChannelOpened {
		status = "up"
	} else {
		status = "down"
	}

	return &spec.GetAgentStatusResponse{
		Status: status,
	}
}
