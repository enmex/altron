package server

import (
	"altron/common/dto"
	"altron/plugin/generated/spec"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s Server) PostPluginsProcess(c *gin.Context) {
	var request spec.ProcessSessionRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /plugins/process json decode failed: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.ProcessSession(&request)
	if err != nil {
		s.log.Errorf("POST /plugins/process: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetPlugins(c *gin.Context) {
	res := s.useCase.GetAllPlugins()
	c.AbortWithStatusJSON(http.StatusOK, res)
}
