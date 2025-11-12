package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tracer/delivery"
	"tracer/pkg/logging"
	"tracer/pkg/tracer"

	"go.uber.org/zap"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	log := logging.NewLogger(os.Getenv("LOG_LEVEL"), os.Getenv("ENVIRONMENT"))
	ctx = logging.WithLogger(ctx, log)
	defer func() {
		done()
		if r := recover(); r != nil {
			log.Errorw("worker went wrong. Panic", zap.Any("panic", r))
		}
	}()

	err := run(ctx)
	if err != nil {
		log.Errorw("failed to run", zap.Error(err))
	}
}

func run(ctx context.Context) error {
	log := logging.FromContext(ctx)
	tracerShutdown := tracer.InitProvider()
	defer tracerShutdown()

	httpHandler := delivery.NewHTTPHandler()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler.Routes(ctx),
	}

	errCh := make(chan error, 1)
	go func() {
		log.Infow("HTTP server starting", "addr", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("HTTP server stopped unexpectedly", zap.Error(err))
		}
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		log.Infow("shutting down HTTP server due to context cancellation")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("failed to shutdown HTTP server gracefully", zap.Error(err))
			return err
		}
		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		log.Infow("HTTP server shutdown complete")
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("failed to serve HTTP server", zap.Error(err))
			return err
		}
		log.Infow("HTTP server stopped")
		return nil
	}
}
