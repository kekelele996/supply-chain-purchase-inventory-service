package router

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/middleware"
)

// registerPurchaseRoutes 采购单模块路由。
func (r *Router) registerPurchaseRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	purchases := api.Group("/purchases", authMW)
	purchases.GET("", r.purchaseHandler.List)
	purchases.GET("/", r.purchaseHandler.List)
	purchases.POST("", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Create)
	purchases.POST("/", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Create)
	purchases.GET("/:id", r.purchaseHandler.Get)
	purchases.PUT("/:id", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Update)
	purchases.DELETE("/:id", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Delete)
	purchases.PUT("/:id/submit", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Submit)
	purchases.PUT("/:id/approve", middleware.RequireRoles(constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Approve)
	purchases.PUT("/:id/reject", middleware.RequireRoles(constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Reject)
	purchases.PUT("/:id/complete", middleware.RequireRoles(constants.RoleOperator, constants.RoleManager, constants.RoleAdmin), r.purchaseHandler.Complete)
}
