package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UseErrorHandling(r *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"error": recovered}})
	}))
	r.Use(newErrorHandler())
}

func newErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// status -1 doesn't overwrite existing status code
			c.JSON(-1, gin.H{"errors": c.Errors})
		}
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// status -1 doesn't overwrite existing status code
			c.JSON(-1, gin.H{"errors": c.Errors})
		}
	}
}
