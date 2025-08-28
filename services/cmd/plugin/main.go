package main

import (
	"altron/config"
	httpServer "altron/pkg/http"
	"altron/pkg/logger"
	"altron/pkg/plugin"
	"altron/plugin/generated/spec"
	"altron/plugin/interfaces"
	"altron/plugin/server"
	"altron/plugin/usecases"
	"strconv"

	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	log := &logger.Logger
	pluginCfg := config.NewPluginConfig()
	appCfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	//plugin manager
	pluginManager, err := plugin.NewPluginManager[interfaces.PluginInterface](pluginCfg)
	if err != nil {
		log.Fatalln(err)
	}

	//usecase
	useCase := usecases.NewPluginUseCase(log, pluginManager)

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
	httpSrv := httpServer.NewServer(ctx, strconv.Itoa(appCfg.AltronPluginPort), handler)

	log.Infoln("Altron plugin server has been started")
	err = httpServer.StartServer(httpSrv)
	if err != nil {
		log.Fatalln(err)
	}
}
