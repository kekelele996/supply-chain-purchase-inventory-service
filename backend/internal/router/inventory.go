package router

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/middleware"
)

// registerInventoryRoutes 库存模块路由。
func (r *Router) registerInventoryRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	inventory := api.Group("/inventory", authMW)
	inventory.GET("", r.inventoryHandler.List)
	inventory.GET("/", r.inventoryHandler.List)
	inventory.POST("", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.inventoryHandler.Create)
	inventory.POST("/", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.inventoryHandler.Create)
	inventory.GET("/alerts", r.inventoryHandler.Alerts)
	inventory.GET("/:id", r.inventoryHandler.Get)
	inventory.PUT("/:id", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.inventoryHandler.Update)
	inventory.PUT("/:id/quantity", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.inventoryHandler.AdjustQuantity)
	inventory.DELETE("/:id", middleware.RequireRoles(constants.RoleAdmin), r.inventoryHandler.Delete)
}
