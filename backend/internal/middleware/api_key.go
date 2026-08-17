package middleware

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

const apiKeyHeader = "X-API-Key"

// VerifyAPIKey 校验 X-API-Key 是否与配置的 API_KEY_SECRET 一致（常量时间比较）。
// 匹配时返回认证用户（服务账号角色 admin），不匹配返回 nil。
func VerifyAPIKey(c *gin.Context, secret string) *util.AuthUser {
	key := c.GetHeader(apiKeyHeader)
	if key == "" {
		return nil
	}
	if subtle.ConstantTimeCompare([]byte(key), []byte(secret)) != 1 {
		util.L(c).Warn("api key mismatch", "request_id", util.GinRequestID(c))
		return nil
	}
	return &util.AuthUser{ID: 0, Username: "api_key_service", Role: constants.RoleAdmin, FromAPIKey: true}
}
