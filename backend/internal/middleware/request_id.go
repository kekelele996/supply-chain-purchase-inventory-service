// Package middleware 提供认证、RBAC、限流、审计等横切中间件。
package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/util"
)

// RequestID 为每个请求生成或透传 X-Request-ID。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = util.NewRequestID()
		}
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Request = c.Request.WithContext(util.WithRequestID(c.Request.Context(), reqID))
		c.Next()
	}
}
