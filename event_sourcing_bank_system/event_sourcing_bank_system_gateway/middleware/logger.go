package middleware

import (
	"event_sourcing_bank_system_gateway/package/contxt"
	"event_sourcing_bank_system_gateway/package/logger"

	"github.com/gin-gonic/gin"
)

func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		log := logger.FromContext(ctx)
		if reqID := contxt.RequestIDFromCtx(ctx); reqID != "" {
			log = log.With("request_id", reqID)
		}

		ctx = logger.WithLogger(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
