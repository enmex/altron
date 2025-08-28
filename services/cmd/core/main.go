package main

import (
	commonHandlers "altron/common/handlers"
	common "altron/common/models"
	"altron/config"
	"altron/core/generated/spec"
	"altron/core/models"
	"altron/core/repositories"
	"altron/core/server"
	"altron/core/usecases"
	"altron/pkg/auth"
	"altron/pkg/db"
	httpServer "altron/pkg/http"
	"altron/pkg/logger"
	"altron/pkg/redis"
	"altron/pkg/request"
	"altron/pkg/sftp"
	"altron/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	logger.Init()

	log := &logger.Logger

	appCfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatalln(err)
	}

	dbCfg, err := config.NewDBConfig()
	if err != nil {
		log.Fatalln(err)
	}

	amqpCfg, err := config.NewAMQPConfig()
	if err != nil {
		log.Fatalln(err)
	}

	redisCfg, err := config.NewRedisConfig()
	if err != nil {
		log.Fatalln(err)
	}

	cfg := config.NewConfig(
		appCfg,
		config.NewAuthConfig(),
		dbCfg,
		config.NewPluginConfig(),
		amqpCfg,
		config.NewSFTPConfig(),
		redisCfg,
	)
	ctx := context.Background()

	// Database
	if err := db.Migrations(cfg.DB, log); err != nil {
		log.Warnln("Unable to migrate database. Retrying in 5s...")
		time.Sleep(5 * time.Second)
		if err := db.Migrations(cfg.DB, log); err != nil {
			log.Fatalln(err)
		}
	}

	dbClient, err := server.NewDBClient(ctx, cfg.DB)
	if err != nil {
		log.Warnln("Unable to get database driver. Retrying in 5s...")
		time.Sleep(5 * time.Second)
		dbClient, err = server.NewDBClient(ctx, cfg.DB)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Infoln("Database connection established, schema migrated successfully")
	defer dbClient.Close()

	// repository
	repo := repositories.NewRepository(dbClient)

	oapi, err := server.NewOpenAPI()
	if err != nil {
		log.Fatalln(err)
	}

	//auth jwt
	jwtAuth := auth.NewJwtAuthenticate(cfg.Auth)
	auth := auth.NewAuthenticate(oapi)

	redisClient := redis.NewRedisCache[[]spec.Characteristic](redisCfg)
	log.Infoln("Redis connection established")

	authUseCase := usecases.NewAuthUseCase(cfg, &jwtAuth, repo.User, repo.Workspace)

	//sign ip manager
	authUseCase.SignUpManager(ctx)

	//sign up admin
	user, err := authUseCase.SignUp(ctx, &spec.AuthRequest{
		Password: cfg.SFTP.Password,
	})
	if err != nil {
		if !errors.Is(err, models.ErrorUserAlreadyExists) {
			log.Fatalln(err)
		}
		user, err = authUseCase.AuthUser(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}

	userInfo, err := authUseCase.GetUserInfo(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	decryptedPassword, err := utils.DecryptString(userInfo.Password, cfg.Auth.HashSalt)
	if err != nil {
		log.Fatalln(err)
	}

	cfg.SFTP.Password = decryptedPassword
	log.Infoln(cfg.SFTP)
	sftpClient, err := sftp.NewClient(cfg.SFTP)
	if err != nil {
		log.Fatalln(err)
	}

	serviceUseCase := usecases.NewServiceUseCase(cfg, redisClient, repo.Service)
	dashboardUseCase := usecases.NewDashboardUseCase(repo.Service, repo.PcapWorkspace)

	//adding analyzer cache to redis
	servicesRes, err := dashboardUseCase.GetDashboard(ctx, uuid.MustParse(user.Id))
	if err != nil {
		log.Fatalln(err)
	}
	characteristicTypes := []string{"ttl", "requests", "ua", "timestamps"}
	for _, service := range servicesRes.Services {
		for _, componentName := range characteristicTypes {
			if err := redisClient.Set(ctx, fmt.Sprintf("%d_%s", service.Port, componentName), make([]spec.Characteristic, 0)); err != nil {
				log.Fatalln(err)
			}
		}
	}

	//load plugins
	res, err := request.Get[spec.GetAllPluginsResponse](fmt.Sprintf("http://%s:%d/plugins", cfg.App.AltronHost, cfg.App.AltronPluginPort))
	if err != nil {
		log.Warnln("unable to load plugins, retrying in 10s...")
		time.Sleep(10 * time.Second)
		res, err = request.Get[spec.GetAllPluginsResponse](fmt.Sprintf("http://%s:%d/plugins", cfg.App.AltronHost, cfg.App.AltronPluginPort))
		if err != nil {
			log.Fatalln(err)
		}
	}

	//load plugins to db
	pluginUseCase := usecases.NewPluginUseCase(repo.Plugin)
	if err := pluginUseCase.CreatePlugins(ctx, res.Plugins); err != nil {
		log.Fatalln(err)
	}
	log.Infof("Plugins %s has been loaded", res.Plugins)

	//load analyzer characteristics to db
	if err := repo.AnalyzerPayload.CreateAnalyzerComponents(ctx, common.ComponentNames); err != nil {
		log.Fatalln(err)
	}

	analyzerHandler := commonHandlers.NewAnalyzerHandler(appCfg, nil)

	//useCase
	useCase := usecases.NewUseCase(
		usecases.NewAdminUseCase(cfg, repo.User, repo.Admin),
		authUseCase,
		dashboardUseCase,
		serviceUseCase,
		usecases.NewWorkspaceUseCase(cfg.App, redisClient, repo.Workspace, repo.Session),
		usecases.NewSessionUseCase(log, repo.Session, repo.Cart, analyzerHandler),
		usecases.NewConversionUseCase(cfg.App, sftpClient),
		pluginUseCase,
		usecases.NewFilterUseCase(log, repo.Filter, repo.Session),
		usecases.NewSessionAnalyzerUseCase(repo.AnalyzerPayload, repo.Workspace, redisClient),
		usecases.NewCartUseCase(repo.Cart),
		usecases.NewPcapWorkspaceUseCase(cfg.App, sftpClient, repo.PcapWorkspace),
	)

	// Server
	srv := server.NewServer(
		ctx,
		log,
		cfg,
		useCase,
	)

	options, err := server.NewServerOptions(cfg, auth, &jwtAuth, repo.User)
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	handler := spec.RegisterHandlersWithOptions(r, srv, options)
	httpSrv := httpServer.NewServer(ctx, strconv.Itoa(cfg.App.AltronPort), handler)

	log.Infoln("Altron server has been started")
	err = httpServer.StartServer(httpSrv)
	if err != nil {
		log.Fatalln(err)
	}
}
