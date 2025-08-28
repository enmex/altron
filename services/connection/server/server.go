package server

import (
	"context"
	"net/http"

	"altron/connection/generated/spec"
	"altron/connection/usecases"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	context  context.Context
	log      *logrus.Logger
	useCase  *usecases.UseCase
	upgrader websocket.Upgrader
}
type OpenAPI struct {
	*openapi3.T
	routers.Router
}

func NewServer(
	context context.Context,
	logger *logrus.Logger,
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
		log:      logger,
		useCase:  useCase,
		upgrader: upgrader,
	}
}

func NewServerOptions() (spec.GinServerOptions, error) {
	return spec.GinServerOptions{
		Middlewares: []spec.MiddlewareFunc{},
	}, nil
}

func NewOpenAPI() (*OpenAPI, error) {
	var api OpenAPI
	var err error

	api.T, err = spec.GetSwagger()
	if err != nil {
		return nil, err
	}
	api.Servers = openapi3.Servers{&openapi3.Server{URL: "/api/v1"}, &openapi3.Server{URL: "/"}}
	api.Router, err = legacy.NewRouter(api.T)
	if err != nil {
		return nil, err
	}
	return &api, nil
}
