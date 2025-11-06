package main

import (
	"clean_architect/config"
	"clean_architect/env"
	"clean_architect/package/server"
	"clean_architect/presentation"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		done()
		if r := recover(); r != nil {
			log.Printf("ERROR: worker went wrong. Panic: %v", r)
		}
	}()

	err := realMain(ctx)
	done()
	if err != nil {
		log.Printf("ERROR: realMain has failed: %v", err)
		return
	}
	log.Println("INFO: Worker shutdown successful")

}

func realMain(ctx context.Context) error {
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Printf("ERROR: LoadConfig has failed: %v", err)
		return err
	}
	env, err := env.LoadEnv(ctx, cfg)
	if err != nil {
		log.Printf("ERROR: LoadEnv has failed: %v", err)
		return err
	}
	presenter := presentation.NewPresenter(cfg, env)
	srv, err := server.New(cfg.Server.Port)
	if err != nil {
		return err
	}
	go func() {
		err = srv.ServeHTTPHandler(ctx, presenter.AppHttp().Routes(ctx))
		if err != nil {
			log.Printf("ERROR: Serve HTTP Handler failed err=%v", err)
		}
	}()

	// wait for signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	log.Printf("INFO: Receive os.Signal: %s", s.String())
	log.Println("INFO: Shutting down worker ...")
	return nil
}
