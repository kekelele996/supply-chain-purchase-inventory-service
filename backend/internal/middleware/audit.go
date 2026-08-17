package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

// Audit 请求审计中间件：记录 request_id/method/path/status/latency/user。
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := util.GinRequestID(c)
		util.L(c).Info(constants.LogTplRequestIn, "method", c.Request.Method, "path", c.Request.URL.Path, "ip", c.ClientIP())
		c.Next()
		user := util.CurrentUser(c)
		util.L(c).Info(constants.LogTplRequestDone,
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"user_id", user.ID,
		)
	}
}
