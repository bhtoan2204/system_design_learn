package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type LimitCfg struct {
	Max    int
	Window time.Duration
}

// perRoute is a map of route to limit configuration
func RateLimitSlidingWindow(
	redisClient *redis.Client,
	defaultCfg LimitCfg,
	perRoute map[string]LimitCfg,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}
		method := c.Request.Method
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		// Skip rate limiting for health-check endpoints
		if route == "/health-check" {
			c.Next()
			return
		}

		cfg := defaultCfg
		if perRoute != nil {
			if rCfg, ok := perRoute[method+" "+route]; ok {
				cfg = rCfg
			}
		}

		ip := c.ClientIP()
		key := buildKey(ip, method, route)

		if allow(redisClient, c.Request.Context(), key, cfg.Max, cfg.Window) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "rate limit exceeded",
		})
	}
}

func buildKey(ip, method, route string) string {
	return fmt.Sprintf("rl:%s:%s:%s", ip, method, route)
}

func allow(rdb *redis.Client, ctx context.Context, key string, maxRequests int, window time.Duration) bool {
	nowMs := time.Now().UnixMilli()
	startMs := nowMs - window.Milliseconds()

	member := strconv.FormatInt(nowMs, 10) + ":" + strconv.Itoa(rand.Int())

	pipe := rdb.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", startMs))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(nowMs), Member: member})
	countCmd := pipe.ZCount(ctx, key, fmt.Sprintf("%d", startMs), "+inf")
	pipe.Expire(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return true
	}
	n, _ := countCmd.Result()
	return n <= int64(maxRequests)
}
