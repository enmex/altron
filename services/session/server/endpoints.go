package server

import (
	commonDto "altron/common/dto"
	"altron/session/generated/spec"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s Server) PostPorts(c *gin.Context) {
	ctx := c.Request.Context()

	var request spec.CreatePortsRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /ports json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Packet.AddPcapPorts(ctx, &request); err != nil {
		s.log.Errorf("POST /ports: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

func (s Server) DeletePorts(c *gin.Context) {
	ctx := c.Request.Context()

	var request spec.DeletePortRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("DELETE /ports json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Packet.DeletePcapPort(ctx, &request); err != nil {
		s.log.Errorf("DELETE /ports: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) PatchLogs(c *gin.Context) {
	ctx := c.Request.Context()

	var request spec.UpdateContainerRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PATCH /logs json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.LogsCollector.UpdateContainer(ctx, &request); err != nil {
		s.log.Errorf("PATCH /logs: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) PostPcaps(c *gin.Context) {
	ctx := c.Request.Context()

	var request spec.UploadPcapRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /pcaps: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Packet.ImportFile(ctx, &request); err != nil {
		s.log.Errorf("POST /pcaps: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) GetLogsContainerID(c *gin.Context, containerID spec.ContainerID) {
	res, err := s.useCase.LogsCollector.GetCachedLogs(string(containerID))
	if err != nil {
		s.log.Errorf("GET /logs/%v: %v", containerID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}
