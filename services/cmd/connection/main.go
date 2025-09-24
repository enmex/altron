package main

import (
	commonHandlers "altron/common/handlers"
	common "altron/common/models"
	"altron/config"
	"altron/connection/generated/spec"
	"altron/connection/handlers"
	"altron/connection/server"
	"altron/connection/usecases"
	"altron/pkg/amqp"
	httpServer "altron/pkg/http"
	"altron/pkg/logger"
	"altron/pkg/redis"
	"altron/pkg/request"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	log := &logger.Logger

	appCfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatalln(err)
	}
	redisCfg, err := config.NewRedisConfig()
	if err != nil {
		log.Fatalln(err)
	}
	amqpCfg, err := config.NewAMQPConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	//message queue client
	mqClient, err := amqp.NewClient(amqpCfg)
	if err != nil {
		for i := 0; i < 3 && err != nil; i++ {
			log.Warnln("Unable to connect to RabbitMQ. Retrying in 10s...")
			time.Sleep(10 * time.Second) //epic hack
			mqClient, err = amqp.NewClient(amqpCfg)
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Infoln("AMQP connection established")

	redisClient := redis.NewRedisCache[[]common.Characteristic](redisCfg)
	log.Infoln("Redis connection established")

	sessionHandler := handlers.NewSessionHandler(log, appCfg)
	analyzerHandler := commonHandlers.NewAnalyzerHandler(appCfg, redisClient)

	pcapWorkspaceUsecase, err := usecases.NewPcapWorkspaceUseCase(log, appCfg, mqClient)
	if err != nil {
		log.Fatalln(err)
	}

	useCase := usecases.NewUseCase(
		usecases.NewConnectionUseCase(log, appCfg, mqClient, sessionHandler),
		usecases.NewWorkspaceUseCase(log, appCfg, mqClient, sessionHandler, analyzerHandler),
		pcapWorkspaceUsecase,
	)

	// reset all workspaces
	if _, err := request.PostWithEmptyResponse(
		fmt.Sprintf("http://altron.core.loc:%d/api/workspaces/reset", appCfg.AltronPort),
		nil,
	); err != nil {
		time.Sleep(10 * time.Second)
		if _, err := request.PostWithEmptyResponse(
			fmt.Sprintf("http://altron.core.loc:%d/api/workspaces/reset", appCfg.AltronPort),
			nil,
		); err != nil {
			log.Fatalln(err)
		}
	}

	// Server
	srv := server.NewServer(
		ctx,
		log,
		useCase,
	)

	options, err := server.NewServerOptions()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	handler := spec.RegisterHandlersWithOptions(r, srv, options)
	httpSrv := httpServer.NewServer(ctx, strconv.Itoa(appCfg.AltronConnectionPort), handler)

	log.Infoln("Altron connection server has been started")
	err = httpServer.StartServer(httpSrv)
	if err != nil {
		log.Fatalln(err)
	}
}
