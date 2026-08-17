package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

// AuthConfig 认证中间件配置。
type AuthConfig struct {
	JWTSecret    string
	APIKeySecret string
}

// Auth JWT + API Key 双认证中间件：优先校验 X-API-Key，其次校验 Authorization Bearer JWT。
func Auth(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiUser := VerifyAPIKey(c, cfg.APIKeySecret); apiUser != nil {
			util.SetAuthUser(c, *apiUser)
			util.L(c).Info("api key authenticated", "user", apiUser.Username)
			c.Next()
			return
		}
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MessageOf(constants.CodeUnauthorized))
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := util.ParseAccessToken(cfg.JWTSecret, tokenStr)
		if err != nil {
			code := constants.CodeTokenInvalid
			if err == util.ErrExpiredToken {
				code = constants.CodeTokenExpired
			}
			util.Fail(c, http.StatusUnauthorized, code, constants.MessageOf(code))
			return
		}
		util.SetAuthUser(c, util.AuthUser{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     constants.UserRole(claims.Role),
		})
		c.Next()
	}
}
