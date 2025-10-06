package main

import (
	"context"
	"event_sourcing_bank_system_api/package/logger"
	"event_sourcing_bank_system_api/presentation"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := logger.DefaultLogger()
	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Error("application went wrong. Panic err=%v", zap.Error(r.(error)))
		}
	}()
	start(ctx)
}

func start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	app, err := presentation.NewApp(ctx)
	if err != nil {
		log.Error("NewApp failed", zap.Error(err))
		return fmt.Errorf("new app got err=%w", err)
	}

	if app == nil {
		log.Error("NewApp returned nil app without error")
		return fmt.Errorf("NewApp returned nil app without error")
	}
	log.Info("Starting application")
	return app.Start(ctx)
}
