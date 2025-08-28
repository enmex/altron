package server

import (
	commonDto "altron/common/dto"
	"altron/core/generated/spec"
	"altron/core/middleware"
	"altron/core/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s Server) PostAuthSignIn(c *gin.Context) {
	ctx := c.Request.Context()
	var request spec.AuthRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /auth/signIn json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
			Code:    commonDto.ErrorBadRequest,
			Message: err.Error(),
		})
		return
	}
	res, err := s.useCase.Auth.SignIn(ctx, &request)
	if err != nil {
		s.log.Errorf("POST /auth/signIn: %v", err)
		if errors.Is(err, models.ErrorInvalidCredentials) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, res)
}

func (s Server) GetDashboard(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.Dashboard.GetDashboard(ctx, userID)
	if err != nil {
		s.log.Errorf("GET /dashboard: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostServices(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request *spec.CreateServiceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /services json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
			Code:    commonDto.ErrorBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Service.CreateService(ctx, userID, request)
	if err != nil {
		s.log.Errorf("POST /services: %v", err)
		if errors.Is(err, models.ErrorServiceAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, res)
}

func (s Server) PatchServicesServiceID(c *gin.Context, serviceID spec.ServiceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.UpdateServiceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PATCH /services/%v json decode failed: %v", serviceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Service.UpdateService(ctx, userID, uuid.MustParse(string(serviceID)), &request); err != nil {
		s.log.Errorf("PATCH /services/%v: %v", serviceID, err)
		if errors.Is(err, models.ErrorAllFieldsEmpty) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) DeleteServicesServiceID(c *gin.Context, serviceID spec.ServiceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	if err := s.useCase.Service.DeleteService(ctx, userID, uuid.MustParse(string(serviceID))); err != nil {
		s.log.Errorf("DELETE /services/%v: %v", serviceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) PostWorkspaces(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.CreateWorkspaceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /workspaces json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Workspace.CreateWorkspace(ctx, userID, &request)
	if err != nil {
		s.log.Errorf("POST /workspaces: %v", err)
		if errors.Is(err, models.ErrorWorkspaceAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, res)
}

func (s Server) GetServicesServicePort(c *gin.Context, servicePort spec.ServicePort) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.Service.GetService(ctx, userID, uint16(servicePort))
	if err != nil {
		s.log.Errorf("GET /services/%v: %v", servicePort, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostWorkspacesWorkspaceIDSessions(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.AddSessionsToWorkspaceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /workspaces/%v/sessions json decode failed: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Session.AddSessions(ctx, uuid.MustParse(string(workspaceID)), &request); err != nil {
		s.log.Errorf("POST /workspaces/%v/sessions: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) DeleteWorkspacesWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	if err := s.useCase.Workspace.DeleteWorkspace(ctx, userID, uuid.MustParse(string(workspaceID))); err != nil {
		s.log.Errorf("DELETE /workspaces/%v: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) PatchWorkspacesWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.UpdateWorkspaceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PATCH /workspaces/%v json decode failed: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Workspace.UpdateWorkspace(ctx, userID, uuid.MustParse(string(workspaceID)), &request); err != nil {
		s.log.Errorf("PATCH /workspaces/%v: %v", workspaceID, err)
		if errors.Is(err, models.ErrorWorkspaceNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) PostConversionsSessions(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.ConvertSessionToExploitRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /sessions/conversion json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	res, err := s.useCase.Conversion.ConvertSessionToExploit(ctx, &request)
	if err != nil {
		s.log.Errorf("POST /sessions/conversion: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) DeleteWorkspacesWorkspaceIDSessions(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	if err := s.useCase.Session.ClearSessions(ctx, uuid.MustParse(string(workspaceID))); err != nil {
		s.log.Errorf("DELETE /workspaces/%v/sessions: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) GetPlugins(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Plugin.GetAllPlugins(ctx)
	if err != nil {
		s.log.Errorf("GET /plugins: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetServicesScan(c *gin.Context, params spec.GetServicesScanParams) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	scope := ""
	if params.Scope != nil {
		scope = string(*params.Scope)
	}

	res, err := s.useCase.Service.ScanHost(ctx, scope)
	if err != nil {
		s.log.Errorf("GET /services/scan: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetServicesScanServicePort(c *gin.Context, servicePort spec.ServicePort) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Service.ScanService(ctx, uint16(servicePort))
	if err != nil {
		s.log.Errorf("GET /services/scan/%v: %v", servicePort, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetWorkspacesWorkspaceIDSessions(c *gin.Context, workspaceID spec.WorkspaceID, params spec.GetWorkspacesWorkspaceIDSessionsParams) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Session.GetPaginatedSessions(ctx, uuid.MustParse(string(workspaceID)), int(params.Pagination))
	if err != nil {
		s.log.Errorf("GET /workspaces/%v/sessions/?pagination=%v: %v", workspaceID, params.Pagination, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetFilters(c *gin.Context, params spec.GetFiltersParams) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.Filter.GetAllFilters(ctx, userID, params.ServicePort)
	if err != nil {
		s.log.Errorf("GET /filters: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostFilters(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.CreateFilterRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /filters json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Filter.CreateFilter(ctx, userID, &request); err != nil {
		s.log.Errorf("POST /filters: %v", err)
		if errors.Is(err, models.ErrorDuplicateFilter) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusCreated)
}

func (s Server) DeleteFiltersFilterID(c *gin.Context, filterID spec.FilterID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	if err := s.useCase.Filter.DeleteFilter(ctx, userID, uuid.MustParse(string(filterID))); err != nil {
		s.log.Errorf("DELETE /filters/%v: %v", filterID, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) PatchFiltersFilterID(c *gin.Context, filterID spec.FilterID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.UpdateFilterRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PATCH /filters/%v json decode failed: %v", filterID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Filter.UpdateFilter(ctx, userID, &request, uuid.MustParse(string(filterID))); err != nil {
		s.log.Errorf("PATCH /filters/%v: %v", filterID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) PostConversionsPackets(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.ConvertPacketToExploitRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /conversion/packets json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Conversion.ConvertPacketToExploit(ctx, &request)
	if err != nil {
		s.log.Errorf("POST /conversion/packets: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostWorkspacesWorkspaceIDSessionsSearch(c *gin.Context, workspaceID spec.WorkspaceID, params spec.PostWorkspacesWorkspaceIDSessionsSearchParams) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.SearchSessionsRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("GET /workspaces/%v/sessions/search?pagination=%v json decode failed: %v", workspaceID, params.Pagination, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Session.SearchSessions(ctx, uuid.MustParse(string(workspaceID)), int(params.Pagination), &request)
	if err != nil {
		s.log.Errorf("GET /workspaces/%v/sessions/search?pagination=%v: %v", workspaceID, params.Pagination, err)
		if errors.Is(err, models.ErrorInvalidSearchInput) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostSessions(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.CreateSessionRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /sessions json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Session.CreateSession(ctx, userID, &request); err != nil {
		s.log.Errorf("POST /sessions: %v", err)
		if errors.Is(err, models.ErrorDuplicateSession) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

func (s Server) GetSessionsSessionID(c *gin.Context, sessionID spec.SessionID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Session.GetSession(ctx, uuid.MustParse(string(sessionID)))
	if err != nil {
		s.log.Errorf("GET /sessions/%v: %v", sessionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostSessionAnalyzerWorkspacesWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.CreateAnalyzerPayloadRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /session-analyzer/workspaces/%v json decode failed: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.SessionAnalyzer.CreateAnalyzerPayload(ctx, userID, uuid.MustParse(string(workspaceID)), &request); err != nil {
		s.log.Errorf("POST /session-analyzer/workspaces/%v: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

func (s Server) GetSessionAnalyzerComponents(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.SessionAnalyzer.GetAllAnalyzerComponents(ctx)
	if err != nil {
		s.log.Errorf("POST /session-analyzer/components: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetSessionAnalyzerServicesServicePort(c *gin.Context, servicePort spec.ServicePort) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.SessionAnalyzer.GetAnalyzerPayload(ctx, userID, uint16(servicePort))
	if err != nil {
		s.log.Errorf("GET /session-analyzer/%v: %v", servicePort, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetSessionAnalyzerWorkspacesWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.SessionAnalyzer.GetWorkspaceAnalyzerPayload(ctx, userID, uuid.MustParse(string(workspaceID)))
	if err != nil {
		s.log.Errorf("GET /session-analyzer/%v: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetSessionAnalyzerServicesServicePortChecker(c *gin.Context, servicePort spec.ServicePort) {
	if c.IsAborted() {
		s.log.Infoln("Hereeeeeee")
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.SessionAnalyzer.GetServiceCheckerMask(ctx, userID, uint16(servicePort))
	if err != nil {
		s.log.Errorf("GET /session-analyzer/services/%v/checker: %v", servicePort, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostCartsWorkspaceIDSessionsMerge(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	res, err := s.useCase.Session.MergeSessions(ctx, userID, uuid.MustParse(string(workspaceID)))
	if err != nil {
		s.log.Errorf("POST /workspaces/%v/merge: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) DeleteCartsWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	if err := s.useCase.Cart.ClearCart(ctx, uuid.MustParse(string(workspaceID))); err != nil {
		s.log.Errorf("DELETE /workspaces/%v/baskets: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) GetCartsWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID, params spec.GetCartsWorkspaceIDParams) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Cart.GetCart(ctx, uuid.MustParse(string(workspaceID)), int(params.Pagination))
	if err != nil {
		s.log.Errorf("GET /workspaces/%v/baskets: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostCartsWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.AddSessionsToCartRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /workspaces/%v/baskets json decode failed: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Cart.AddSessions(ctx, uuid.MustParse(string(workspaceID)), &request); err != nil {
		s.log.Errorf("POST /workspaces/%v/baskets: %v", workspaceID, err)
		if errors.Is(err, models.ErrorDuplicateSession) {
			c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
				Code:    commonDto.ErrorBadRequest,
				Message: err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) DeleteCartsWorkspaceIDSessionID(c *gin.Context, workspaceID spec.WorkspaceID, sessionID spec.SessionID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	if err := s.useCase.Cart.RemoveSession(ctx, uuid.MustParse(string(workspaceID)), uuid.MustParse(string(sessionID))); err != nil {
		s.log.Errorf("DELETE /workspaces/%v/baskets/%v: %v", workspaceID, sessionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func (s Server) PostConversionsFiles(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.ExtractFilesFromPacketRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /conversion/files json decode failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Conversion.ExtractFiles(ctx, &request)
	if err != nil {
		s.log.Errorf("POST /conversion/files: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetPcapWorkspacesPcapWorkspaceIDSessions(c *gin.Context, pcapWorkspaceID spec.PcapWorkspaceID, params spec.GetPcapWorkspacesPcapWorkspaceIDSessionsParams) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Session.GetPaginatedPcapSessions(ctx, uuid.MustParse(string(pcapWorkspaceID)), int(params.Pagination))
	if err != nil {
		s.log.Errorf("GET /pcap-workspaces/%v/sessions: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostPcaps(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	pcapFile, header, err := c.Request.FormFile("pcap")
	if err != nil {
		s.log.Errorf("POST /pcaps: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	defer pcapFile.Close()

	res, err := s.useCase.PcapWorkspace.CreatePcapWorkspace(ctx, userID, pcapFile, header)
	if err != nil {
		s.log.Errorf("POST /pcaps: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, res)
}

func (s Server) PostPcapWorkspacesPcapWorkspaceIDSessions(c *gin.Context, pcapWorkspaceID spec.PcapWorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.AddPcapSessionRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /pcap-workspaces/%v/sessions: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.Session.AddPcapSession(ctx, uuid.MustParse(string(pcapWorkspaceID)), &request); err != nil {
		s.log.Errorf("POST /pcap-workspaces/%v/sessions: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusCreated)
}

func (s Server) PatchPcapWorkspacesPcapWorkspaceID(c *gin.Context, pcapWorkspaceID spec.PcapWorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	var request spec.UpdatePcapWorkspaceRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("PATCH /pcap-workspaces/%v json decode failed: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := s.useCase.PcapWorkspace.UpdatePcapWorkspace(ctx, userID, uuid.MustParse(string(pcapWorkspaceID)), &request); err != nil {
		s.log.Errorf("PATCH /workspaces/%v/pcaps: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) PostAdminReconfigure(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.ReconfigureAltronRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /admin/reconfigure json decode failed %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Admin.ResetAll(ctx)
	if err != nil {
		s.log.Errorf("POST /admin/reconfigure %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetAdminInfo(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	res, err := s.useCase.Admin.GetInfo(ctx)
	if err != nil {
		s.log.Errorf("GET /admin/info %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostAdminSignIn(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()

	var request spec.AuthRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		s.log.Errorf("POST /admin/signIn json decode failed %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	res, err := s.useCase.Auth.SignInManager(ctx, &request)
	if err != nil {
		s.log.Errorf("POST /admin/signIn %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetLogsContainerID(c *gin.Context, containerID spec.ContainerID) {
	if c.IsAborted() {
		return
	}
	res, err := http.Get(
		fmt.Sprintf("http://%s:%d/logs/%v", s.cfg.App.AltronHost, s.cfg.App.AltronSessionPort, containerID),
	)
	if err != nil {
		s.log.Errorf("GET /logs/%v: %v", containerID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		s.log.Errorf("GET /logs/%v: %v", containerID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
	c.Abort()
}

func (s Server) GetWorkspacesWorkspaceID(c *gin.Context, workspaceID spec.WorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	workspaceUuid := uuid.MustParse(string(workspaceID))
	res, err := s.useCase.Workspace.GetWorkspace(ctx, userID, workspaceUuid)
	if err != nil {
		s.log.Errorf("GET /workspaces/%v: %v", workspaceID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) PostWorkspacesReset(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	if err := s.useCase.Workspace.ResetWorkspaces(ctx, userID); err != nil {
		s.log.Errorf("POST /workspaces/reset: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (s Server) GetUsersInfo(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	res, err := s.useCase.Auth.GetUserInfo(c.Request.Context())
	if err != nil {
		s.log.Errorf("GET /users/info: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonDto.ErrorResponse{
			Code:    commonDto.ErrorInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) GetPcapWorkspacesPcapWorkspaceID(c *gin.Context, pcapWorkspaceID spec.PcapWorkspaceID) {
	if c.IsAborted() {
		return
	}
	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	pcapWorkspaceUuid, err := uuid.Parse(string(pcapWorkspaceID))
	if err != nil {
		s.log.Errorf("GET /pcap-workspaces/%v: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
			Code:    commonDto.ErrorBadRequest,
			Message: models.ErrorPcapWorkspaceNotFound.Error(),
		})
		return
	}

	res, err := s.useCase.PcapWorkspace.GetPcapWorkspace(ctx, userID, pcapWorkspaceUuid)
	if err != nil {
		s.log.Errorf("GET /pcap-workspaces/%v: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
			Code:    commonDto.ErrorBadRequest,
			Message: err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

func (s Server) DeletePcapWorkspacesPcapWorkspaceID(c *gin.Context, pcapWorkspaceID spec.PcapWorkspaceID) {
	if c.IsAborted() {
		return
	}

	ctx := c.Request.Context()
	userID := middleware.GetUserIdFromContext(ctx)

	pcapWorkspaceUuid, err := uuid.Parse(string(pcapWorkspaceID))
	if err != nil {
		s.log.Errorf("DELETE /pcap-workspaces/%v: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
			Code:    commonDto.ErrorBadRequest,
			Message: models.ErrorPcapWorkspaceNotFound.Error(),
		})
		return
	}

	if err := s.useCase.PcapWorkspace.DeletePcapWorkspace(ctx, userID, pcapWorkspaceUuid); err != nil {
		s.log.Errorf("DELETE /pcap-workspaces/%v: %v", pcapWorkspaceID, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, commonDto.ErrorResponse{
			Code:    commonDto.ErrorBadRequest,
			Message: models.ErrorPcapWorkspaceNotFound.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
