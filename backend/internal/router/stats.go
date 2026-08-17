package router

import "github.com/gin-gonic/gin"

// registerStatsRoutes 成本统计模块路由（所有登录角色可查看）。
func (r *Router) registerStatsRoutes(api *gin.RouterGroup, authMW gin.HandlerFunc) {
	stats := api.Group("/stats", authMW)
	stats.GET("/cost/by-category", r.statsHandler.ByCategory)
	stats.GET("/cost/by-supplier", r.statsHandler.BySupplier)
	stats.GET("/cost/trend", r.statsHandler.Trend)
	stats.GET("/cost/summary", r.statsHandler.Summary)
}
