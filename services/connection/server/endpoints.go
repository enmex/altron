package server

import (
	commonDto "altron/common/dto"
	"altron/connection/generated/spec"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s Server) GetConnectSessionsServicePort(c *gin.Context, servicePort spec.ServicePort) {
	ctx := c.Request.Context()

	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Errorln(err)
		return
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Time{})
	ws.SetWriteDeadline(time.Time{})

	if err := s.useCase.Connection.ListenSessions(ctx, ws, uint16(servicePort)); err != nil {
		s.log.Errorln(err)
	}
	if err := ws.Close(); err != nil {
		s.log.Errorf("GET /connect/sessions/%v closed with error: %v", servicePort, err)
	}
}

func (s Server) GetConnectLogsContainerID(c *gin.Context, containerID spec.ContainerID) {
	ctx := c.Request.Context()

	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Errorln(err)
		return
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Time{})
	ws.SetWriteDeadline(time.Time{})

	if err := s.useCase.Connection.ListenLogs(ctx, ws, string(containerID)); err != nil {
		s.log.Errorln(err)
	}
	if err := ws.Close(); err != nil {
		s.log.Errorf("GET /connect/logs/%v closed with error: %v", containerID, err)
	}
}

func (s Server) GetConnectAgent(c *gin.Context) {
	ctx := c.Request.Context()
	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Errorln(err)
		return
	}

	ws.SetReadDeadline(time.Time{})
	ws.SetWriteDeadline(time.Time{})

	if err := s.useCase.Connection.OpenMainChannel(ctx, ws); err != nil {
		s.log.Errorln("Agent main channel error:", err)
	}
	if err := ws.Close(); err != nil {
		s.log.Errorln("Main channel closed with error:", err)
	}
}

func (s Server) PutWorkspacesWorkspaceIDSessionsServicePort(c *gin.Context, workspaceID spec.WorkspaceID, servicePort spec.ServicePort) {
	ctx := c.Request.Context()

	var request spec.StartSessionsListeningRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PUT /workspaces/%v/sessions/%v json decode failed: %v", workspaceID, servicePort, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Workspace.StartListeningSessions(ctx, uuid.MustParse(string(workspaceID)), uint16(servicePort), &request); err != nil {
		s.log.Errorf("PUT /workspaces/%v/sessions/%v: %v", workspaceID, servicePort, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (s Server) DeleteWorkspacesWorkspaceIDSessions(c *gin.Context, workspaceID spec.WorkspaceID) {
	s.useCase.Workspace.StopListeningSessions(uuid.MustParse(string(workspaceID)))
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) PutWorkspacesPcapsPcapWorkspaceID(c *gin.Context, pcapWorkspaceID spec.PcapWorkspaceID) {
	var request spec.StartPcapListeningRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PUT /workspaces/pcaps/%v: %v", pcapWorkspaceID, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	go s.useCase.PcapWorkspace.StartListeningPcap(uuid.MustParse(string(pcapWorkspaceID)), &request)
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) PostEvents(c *gin.Context) {
	ctx := c.Request.Context()

	var request spec.CreateEventRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /events json decode failed: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Connection.CreateEvent(ctx, &request)
	if err != nil {
		s.log.Errorf("POST /events: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, res)
}

func (s Server) GetConnectAgentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	c.AbortWithStatusJSON(http.StatusOK, s.useCase.Connection.GetAgentStatus(ctx))
}
