package server

import (
	"altron/config"
	"altron/pkg/auth"
	"altron/pkg/db"
	"context"

	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/middleware"
	"altron/core/repositories/ent"
	"altron/core/repositories/ent/migrate"
	"altron/core/usecases"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/sirupsen/logrus"
)

type Server struct {
	context context.Context
	log     *logrus.Logger
	cfg     *config.Config
	useCase *usecases.UseCase
}
type OpenAPI struct {
	*openapi3.T
	routers.Router
}

func NewServer(
	context context.Context,
	log *logrus.Logger,
	config *config.Config,
	useCase *usecases.UseCase,
) *Server {
	return &Server{
		context: context,
		log:     log,
		cfg:     config,
		useCase: useCase,
	}
}

func NewServerOptions(
	cfg *config.Config,
	authenticator auth.Authenticator,
	jwtAuthenticator auth.JwtAuthenticator,
	userRepo interfaces.UserRepository,
) (spec.GinServerOptions, error) {
	authMiddleware := middleware.NewAuthMiddleware(cfg, authenticator, jwtAuthenticator, userRepo)

	return spec.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []spec.MiddlewareFunc{
			authMiddleware.HandlerFunc(),
		},
	}, nil
}

func NewDBClient(ctx context.Context, cfg *db.Config) (*ent.Client, error) {
	client, err := ent.Open(cfg.DriverName, cfg.DSN())
	if err != nil {
		return nil, err
	}
	if err := client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return nil, err
	}
	return client, nil
}

func NewOpenAPI() (*OpenAPI, error) {
	var api OpenAPI
	var err error

	api.T, err = spec.GetSwagger()
	if err != nil {
		return nil, err
	}
	api.Servers = openapi3.Servers{&openapi3.Server{URL: "/api"}, &openapi3.Server{URL: "/"}}
	api.Router, err = legacy.NewRouter(api.T)
	if err != nil {
		return nil, err
	}
	return &api, nil
}
