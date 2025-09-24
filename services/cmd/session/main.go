package main

import (
	commonHandler "altron/common/handlers"
	common "altron/common/models"
	"altron/config"
	"altron/pkg/amqp"
	"altron/pkg/request"
	"altron/utils"

	//"altron/pkg/docker"
	httpServer "altron/pkg/http"
	"altron/pkg/logger"
	"altron/pkg/packets"
	"altron/pkg/redis"
	sftp "altron/pkg/sftp"
	"altron/session/dto"
	"altron/session/generated/spec"
	"altron/session/metrics"
	"altron/session/server"
	"altron/session/usecases"
	"context"
	"fmt"
	"strings"
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

	amqpCfg, err := config.NewAMQPConfig()
	if err != nil {
		log.Fatalln(err)
	}

	redisCfg, err := config.NewRedisConfig()
	if err != nil {
		log.Fatalln(err)
	}

	authCfg := config.NewAuthConfig()

	userInfoUrl := fmt.Sprintf("http://altron.core.loc:%d/api/users/info", appCfg.AltronPort)
	userRes, err := request.Get[dto.GetUserInfoResponse](userInfoUrl)
	if err != nil {
		time.Sleep(10 * time.Second)
		userRes, err = request.Get[dto.GetUserInfoResponse](userInfoUrl)
		if err != nil {
			log.Fatalln(err)
		}
	}
	descryptedPassword, err := utils.DecryptString(userRes.Password, authCfg.HashSalt)
	if err != nil {
		log.Fatalln(err)
	}

	sftpCfg := config.NewSFTPConfig()
	sftpCfg.Password = descryptedPassword

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

	analyzerHandler := commonHandler.NewAnalyzerHandler(appCfg, redisClient)
	serverMetrics, err := metrics.NewServerMetrics(log, mqClient)
	if err != nil {
		log.Fatalln(err)
	}

	sessionTree, err := usecases.NewSessionCollectorUseCase(log, appCfg, mqClient, analyzerHandler, serverMetrics)
	if err != nil {
		log.Fatalln(err)
	}

	pcapAnalyzer := packets.NewPcapAnalyzer()

	// dockerClient, err := docker.NewClient()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	sftpClient, err := sftp.NewClient(sftpCfg)
	if err != nil {
		log.Warnln("unable to connect to SFTP, retrying in 5s...")
		time.Sleep(5 * time.Second)
		sftpClient, err = sftp.NewClient(sftpCfg)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Infoln("SFTP client connected")
	if err := sftpClient.MakeDir("pcaps"); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			log.Fatalln(err)
		}
	}
	if err := sftpClient.MakeDir("logs"); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			log.Fatalln(err)
		}
	}
	if err := sftpClient.MakeDir("files"); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			log.Fatalln(err)
		}
	}

	logsCollector, err := usecases.NewLogsCollectorUseCase(log, sftpClient, mqClient)
	if err != nil {
		log.Fatalln(err)
	}

	packetUseCase, err := usecases.NewPacketUseCase(
		log,
		appCfg,
		sftpClient,
		sessionTree,
		logsCollector,
		pcapAnalyzer,
		serverMetrics,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// get ports
	portsUrl := fmt.Sprintf("http://altron.core.loc:%d/api/dashboard", appCfg.AltronPort)
	res, err := request.Get[dto.GetPortsResponse](portsUrl)
	if err != nil {
		time.Sleep(10 * time.Second)
		res, err = request.Get[dto.GetPortsResponse](portsUrl)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if len(res.Services) > 0 {
		portRequests := make([]spec.PortRequest, 0, len(res.Services))
		for _, service := range res.Services {
			portRequests = append(portRequests, spec.PortRequest{
				ContainerID: service.ContainerID,
				Port:        service.Port,
			})
		}
		if err := packetUseCase.AddPcapPorts(ctx, &spec.CreatePortsRequest{
			Ports: portRequests,
		}); err != nil {
			log.Fatalln(err)
		}
	}

	// healthUseCase, err := usecases.NewHealthUseCase(appCfg, log, dockerClient, mqClient)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// //health stats
	// go func(ctx context.Context) {
	// 	if err := healthUseCase.LaunchHealthMonitoring(ctx); err != nil {
	// 		log.Errorln("health", err)
	// 	}
	// }(ctx)

	//session tree cleaner
	go func(ctx context.Context) {
		cleaner := sessionTree.StartSessionCleaner(ctx)
		log.Infoln("Starting the session tree cleaner...")
		for {
			err := <-cleaner.Error()
			log.Errorln("cleaner", err)
		}
	}(ctx)
	go func(ctx context.Context) {
		cleaner := packetUseCase.ClearExpiredClients(ctx)
		for {
			err := <-cleaner.Error()
			log.Errorln("last tcp packets cleaner", err)
		}
	}(ctx)
	
	//servers metrics measurement
	go func() {
		metricsScheduler := packetUseCase.StartMetricsMeasure()
		log.Infoln("Starting the session metrics measurement...")
		for {
			err := <-metricsScheduler.Error()
			log.Errorln("metrics ", err)
		}
	}()

	//servers metrics sender
	go func() {
		senderScheduler := packetUseCase.StartSendMetrics()
		log.Infoln("Starting sending session metrics...")
		for {
			err := <-senderScheduler.Error()
			log.Errorln("sender", err)
		}
	}()

	// Здесь стоит памятник ушедшей эпохе альтрона
	// //producer thread for message queue
	// go func() {
	// 	packetsChan := pcapAnalyzer.GetPacketsChan()
	// 	log.Infoln("Starting the packet producing process...")
	// 	for packet := range packetsChan {
	// 		packetUseCase.ProducePacket(ctx, packet, nil)
	// 	}
	// }()

	//sending pcap dumps
	go func() {
		log.Infoln("Starting monitoring pcap dumps...")
		if err := packetUseCase.MonitorDumps(ctx); err != nil {
			log.Errorln(err)
		}
	}()

	//sending logs
	go func() {
		log.Infoln("Starting monitoring logs bulks...")
		if err := logsCollector.MonitorLogsJson(ctx); err != nil {
			log.Errorln(err)
		}
	}()

	useCase := usecases.NewUseCase(
		packetUseCase,
		sessionTree,
		logsCollector,
	)

	// Server
	srv := server.NewServer(
		ctx,
		log,
		appCfg,
		useCase,
	)

	options, err := server.NewServerOptions()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	handler := spec.RegisterHandlersWithOptions(r, srv, options)
	httpSrv := httpServer.NewServer(ctx, fmt.Sprint(appCfg.AltronSessionPort), handler)

	log.Infoln("Altron session server has been started")
	err = httpServer.StartServer(httpSrv)
	if err != nil {
		log.Fatalln(err)
	}
}
