package router

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/middleware"
)

// registerLogRoutes 操作日志模块路由（manager+）。
func (r *Router) registerLogRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	logs := api.Group("/logs", authMW, middleware.RequireRoles(constants.RoleManager, constants.RoleAdmin))
	logs.GET("", r.logHandler.List)
	logs.GET("/", r.logHandler.List)
	logs.GET("/:id", r.logHandler.Get)
}
