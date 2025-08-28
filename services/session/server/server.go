package server

import (
	"altron/config"
	"context"
	"net/http"

	"altron/session/generated/spec"
	"altron/session/usecases"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	context  context.Context
	log      *logrus.Logger
	cfg      *config.AppConfig
	useCase  *usecases.UseCase
	upgrader websocket.Upgrader
}
type OpenAPI struct {
	*openapi3.T
	routers.Router
}

func NewServer(
	context context.Context,
	log *logrus.Logger,
	config *config.AppConfig,
	useCase *usecases.UseCase,
) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &Server{
		context:  context,
		log:      log,
		cfg:      config,
		useCase:  useCase,
		upgrader: upgrader,
	}
}

func NewServerOptions() (spec.GinServerOptions, error) {
	return spec.GinServerOptions{
		Middlewares: []spec.MiddlewareFunc{},
	}, nil
}
