package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	start(ctx, cancel)
}

func start(ctx context.Context, cancel context.CancelFunc) {
	r := gin.New()
	// r.Use(gin.Recovery())

	r.GET("/panic", func(c *gin.Context) { panic("boom") })
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"success": true}) })

	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error:", err)
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
	fmt.Println("Shutting down worker ...")
}
