package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"university_admission_system/internal/app/container"
	"university_admission_system/pkg/config"
	"university_admission_system/pkg/logger"
	httpapi "university_admission_system/presentation/http"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.NewStdLogger()
	c := container.New(cfg, log)

	if err := c.SeedDemoData(ctx); err != nil {
		log.Error("failed to seed data", err, nil)
		os.Exit(1)
	}

	router := httpapi.NewRouter(httpapi.RouterConfig{
		SubmitService:            c.SubmitApplicationService(),
		IssueOfferService:        c.IssueOfferService(),
		AcceptOfferService:       c.AcceptOfferService(),
		ConfirmEnrollmentService: c.ConfirmEnrollmentService(),
		Logger:                   c.Logger(),
		EnableSwagger:            cfg.EnableSwagger,
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Info("HTTP server listening", map[string]interface{}{"port": cfg.HTTPPort})
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server error", err, nil)
		os.Exit(1)
	}
}
