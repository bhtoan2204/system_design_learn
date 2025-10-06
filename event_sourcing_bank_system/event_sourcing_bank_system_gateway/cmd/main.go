package main

import (
	"context"
	"event_sourcing_bank_system_gateway/application"
	"event_sourcing_bank_system_gateway/package/config"
	"event_sourcing_bank_system_gateway/package/logger"
	"event_sourcing_bank_system_gateway/package/settings"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

//	@title			API Gateway
//	@version		1.0
//	@description	REST -> GRPC API Gateway

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

// @host						localhost:8080
// @BasePath					/
// @query.collection.format	multi
func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	config, _ := config.InitLoadConfig()
	log := logger.NewLogger(config.LogConfig.LogLevel, config.LogConfig.FilePath)
	log.Infow("Server mode: ", zap.Any("config", config))
	ctx = logger.WithLogger(ctx, log)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Info("Application went wrong. Panic err:", zap.Error(r.(error)))
		}
	}()
	err := initialize(ctx, config)
	done()
	if err != nil {
		log.Error("Initialize failed", zap.Error(err))
		return
	}

	log.Info("App shutdown successful")
}

func initialize(ctx context.Context, config *settings.Config) error {
	app, err := application.New(config)
	if err != nil {
		return fmt.Errorf("new app got err=%w", err)
	}

	return app.Start(ctx)
}
