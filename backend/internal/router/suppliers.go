package router

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/middleware"
)

// registerSupplierRoutes 供应商模块路由。
func (r *Router) registerSupplierRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	suppliers := api.Group("/suppliers", authMW)
	suppliers.GET("", r.supplierHandler.List)
	suppliers.GET("/", r.supplierHandler.List)
	suppliers.POST("", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.supplierHandler.Create)
	suppliers.POST("/", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.supplierHandler.Create)
	suppliers.GET("/:id", r.supplierHandler.Get)
	suppliers.PUT("/:id", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.supplierHandler.Update)
	suppliers.DELETE("/:id", middleware.RequireRoles(constants.RoleAdmin), r.supplierHandler.Delete)
	suppliers.PUT("/:id/status", middleware.RequireRoles(constants.RoleManager, constants.RoleAdmin), r.supplierHandler.ChangeStatus)
	suppliers.GET("/:id/inventory", r.supplierHandler.Inventory)
	suppliers.GET("/:id/orders", r.supplierHandler.Orders)
}
