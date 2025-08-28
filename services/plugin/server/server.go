package server

import (
	"context"

	"altron/plugin/generated/spec"
	"altron/plugin/interfaces"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/sirupsen/logrus"
)

type Server struct {
	context context.Context
	log     *logrus.Logger
	useCase interfaces.PluginUseCase
}
type OpenAPI struct {
	*openapi3.T
	routers.Router
}

func NewServer(
	context context.Context,
	logger *logrus.Logger,
	useCase interfaces.PluginUseCase,
) *Server {
	return &Server{
		context: context,
		log:     logger,
		useCase: useCase,
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
