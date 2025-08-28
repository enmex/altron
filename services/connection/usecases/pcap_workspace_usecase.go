package usecases

import (
	common "altron/common/models"
	"altron/connection/dto"
	"altron/connection/generated/spec"
	"altron/connection/interfaces"
	"altron/connection/models"
	"encoding/json"
	"fmt"
	"net/http"

	"altron/config"
	"altron/pkg/amqp"
	"altron/pkg/request"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var _ interfaces.PcapWorkspaceUseCase = (*PcapWorkspaceUseCase)(nil)

type PcapWorkspaceUseCase struct {
	log      *logrus.Logger
	cfg      *config.AppConfig
	client   *amqp.Client
	producer *amqp.Producer
}

func NewPcapWorkspaceUseCase(
	log *logrus.Logger,
	cfg *config.AppConfig,
	client *amqp.Client,
) (*PcapWorkspaceUseCase, error) {
	producer, err := amqp.NewProducer(client)
	if err != nil {
		return nil, err
	}
	return &PcapWorkspaceUseCase{
		log:      log,
		cfg:      cfg,
		client:   client,
		producer: producer,
	}, nil
}

func (u *PcapWorkspaceUseCase) StartListeningPcap(pcapWorkspaceID uuid.UUID, req *spec.StartPcapListeningRequest) {
	consumer, err := amqp.NewConsumer(u.client, uuid.New().String())
	if err != nil {
		u.log.Errorln(err)
		return
	}
	defer consumer.Close()
	if err := u.producer.CreateExchange(req.FileName); err != nil {
		u.log.Errorln(err)
		return
	}
	if err := consumer.Bind(req.FileName); err != nil {
		u.log.Errorln(err)
		return
	}

	u.log.Infof("pcap workspace %s started collecting sessions", pcapWorkspaceID.String())
	defer func() {
		u.log.Infof("pcap workspace %s stopped collecting sessions", pcapWorkspaceID.String())
		consumer.Close()
	}()

	messageChan, err := consumer.Messages()
	if err != nil {
		u.log.Errorln(err)
		return
	}
	for {
		select {
		case message := <-messageChan:
			if string(message.Body) == "eof" {
				statusCode, err := request.PatchWithEmptyResponse(
					fmt.Sprintf("http://%s:%d/api/pcap-workspaces/%s", u.cfg.AltronHost, u.cfg.AltronPort, pcapWorkspaceID.String()),
					dto.UpdatePcapWorkspaceStatusRequest{
						Status: common.COMPLETED,
					},
				)
				if err != nil {
					u.log.Errorln(err)
					return
				}
				if statusCode != http.StatusOK {
					u.log.Errorln(models.ErrorAddingSessionsToWorkspace.Error())
					return
				}
				return
			}
			var session common.Session
			if err := json.Unmarshal(message.Body, &session); err != nil {
				u.log.Errorln(err)
				return
			}
			statusCode, err := request.PostWithEmptyResponse(
				fmt.Sprintf("http://%s:%d/api/pcap-workspaces/%s/sessions", u.cfg.AltronHost, u.cfg.AltronPort, pcapWorkspaceID.String()),
				dto.AddPcapSessionRequest{
					Session: &session,
				},
			)
			if err != nil {
				u.log.Errorln(err)
				return
			}

			if statusCode != http.StatusCreated {
				u.log.Errorln(models.ErrorAddingSessionsToWorkspace.Error())
				return
			}
			u.log.Infof("found session in pcap file %s", req.FileName)
		}
	}
}
