package router

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/middleware"
)

// registerUserRoutes 用户管理模块路由（仅 admin）。
func (r *Router) registerUserRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	users := api.Group("/users", authMW, middleware.RequireRoles(constants.RoleAdmin))
	users.GET("", r.userHandler.List)
	users.GET("/", r.userHandler.List)
	users.POST("", r.userHandler.Create)
	users.POST("/", r.userHandler.Create)
	users.GET("/:id", r.userHandler.Get)
	users.PUT("/:id", r.userHandler.Update)
	users.DELETE("/:id", r.userHandler.Delete)
}
