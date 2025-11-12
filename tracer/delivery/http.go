package delivery

import (
	"errors"
	"net/http"
	"tracer/delivery/middleware"
	"tracer/pkg/logging"

	"context"

	"github.com/gin-gonic/gin"
)

type HTTPHandler interface {
	Routes(ctx context.Context) http.Handler
}

type httpHandler struct {
}

func NewHTTPHandler() HTTPHandler {
	return &httpHandler{}
}

func (h *httpHandler) Routes(ctx context.Context) http.Handler {
	r := gin.New()
	r.Use(middleware.SetRequestID())
	r.Use(middleware.TraceRequest())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		slog := logging.NewSpanLogger(ctx)

		slog.Infow("health-check", "Start health check")
		slog.Infow("health-check", "Running step 1")
		Step1(ctx)

		slog.Warnw("health-check", "Step 2 may be slow", errors.New("slow"))
		Step2(ctx)

		slog.Errorw("health-check", "Step 3 failed", errors.New("timeout"))
		Step3(ctx)

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	return r
}

func Step1(c context.Context) {
	logger := logging.FromContext(c)
	logger.Infow("step 1")
}

func Step2(c context.Context) {
	logger := logging.FromContext(c)
	logger.Infow("step 2")
}

func Step3(c context.Context) {
	logger := logging.FromContext(c)
	logger.Infow("step 3")
}
