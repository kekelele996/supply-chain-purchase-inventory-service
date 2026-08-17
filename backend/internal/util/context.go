package util

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
)

type ctxKey string

const (
	ctxKeyRequestID ctxKey = "request_id"
	ctxKeyUser      ctxKey = "auth_user"
)

// AuthUser 认证后的当前用户信息。
type AuthUser struct {
	ID       uint
	Username string
	Role     constants.UserRole
	FromAPIKey bool
}

// NewRequestID 生成 16 字节随机请求 ID（hex 编码）。
func NewRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(b)
}

// WithRequestID 将 request_id 写入 context。
func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, reqID)
}

// RequestIDFromContext 从 context 读取 request_id。
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// GinRequestID 从 gin.Context 读取 request_id。
func GinRequestID(c *gin.Context) string {
	return RequestIDFromContext(c.Request.Context())
}

// SetAuthUser 将认证用户写入 gin context。
func SetAuthUser(c *gin.Context, u AuthUser) {
	c.Set(string(ctxKeyUser), u)
}

// CurrentUser 从 gin context 读取认证用户。
func CurrentUser(c *gin.Context) AuthUser {
	if v, ok := c.Get(string(ctxKeyUser)); ok {
		if u, ok := v.(AuthUser); ok {
			return u
		}
	}
	return AuthUser{}
}
