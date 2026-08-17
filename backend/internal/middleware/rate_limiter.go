package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

// RateLimiter Redis 固定窗口限流中间件。
type RateLimiter struct {
	client    *redis.Client
	perMinute int
}

// NewRateLimiter 构造限流器。
func NewRateLimiter(client *redis.Client, perMinute int) *RateLimiter {
	return &RateLimiter{client: client, perMinute: perMinute}
}

// Handler 返回限流中间件。
func (rl *RateLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.perMinute <= 0 {
			c.Next()
			return
		}
		key := fmt.Sprintf("supplychain:ratelimit:%s", c.ClientIP())
		ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
		defer cancel()
		window := time.Now().Truncate(time.Minute).Unix()
		windowKey := fmt.Sprintf("%s:%d", key, window)
		count, err := rl.client.Incr(ctx, windowKey).Result()
		if err != nil {
			// Redis 异常时放行，避免限流服务故障阻断主流程。
			c.Next()
			return
		}
		if count == 1 {
			_ = rl.client.Expire(ctx, windowKey, 2*time.Minute).Err()
		}
		if count > int64(rl.perMinute) {
			util.L(c).Warn("rate limited", "ip", c.ClientIP(), "limit", rl.perMinute)
			util.Fail(c, http.StatusTooManyRequests, constants.CodeRateLimited, constants.MessageOf(constants.CodeRateLimited))
			return
		}
		c.Next()
	}
}
