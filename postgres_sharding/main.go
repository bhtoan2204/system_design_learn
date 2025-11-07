package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database_sharding/api"
	"database_sharding/persistent"
)

func main() {
	ctx := context.Background()
	db, err := persistent.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if closeErr := persistent.Close(db); closeErr != nil {
			log.Printf("failed to close database: %v", closeErr)
		}
	}()
	if err := db.WithContext(ctx).AutoMigrate(&persistent.Tenant{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	if err := persistent.EnsureTenantDistribution(ctx, db); err != nil {
		log.Fatalf("failed to distribute tenants table: %v", err)
	}
	if err := persistent.SeedTenants(ctx, db); err != nil {
		log.Fatalf("failed to seed tenants: %v", err)
	}

	apiServer := api.NewServer(db)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: apiServer.Handler(),
	}

	go func() {
		log.Printf("HTTP server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown http server: %v", err)
	}
	log.Println("server stopped gracefully")
}
