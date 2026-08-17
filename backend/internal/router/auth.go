package router

import (
	"github.com/gin-gonic/gin"
)

// registerAuthRoutes 认证模块路由（login/refresh 公开，me/password 需登录）。
func (r *Router) registerAuthRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	auth := api.Group("/auth")
	auth.POST("/login", r.authHandler.Login)
	auth.POST("/refresh", r.authHandler.Refresh)
	secured := auth.Group("", authMW)
	secured.GET("/me", r.authHandler.Me)
	secured.PUT("/password", r.authHandler.ChangePassword)
}
