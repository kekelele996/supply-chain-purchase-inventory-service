// Package router 统一注册全部路由。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/config"
	"github.com/supplychain/supplychain-api/internal/handler"
	"github.com/supplychain/supplychain-api/internal/middleware"
	"github.com/supplychain/supplychain-api/internal/service"
)

// Router 持有全部处理器依赖并组装路由。
type Router struct {
	cfg                *config.Config
	authHandler        *handler.AuthHandler
	userHandler        *handler.UserHandler
	supplierHandler    *handler.SupplierHandler
	inventoryHandler   *handler.InventoryHandler
	purchaseHandler    *handler.PurchaseHandler
	statsHandler       *handler.StatsHandler
	logHandler         *handler.LogHandler
	rateLimiter        *middleware.RateLimiter
	operationLogService *service.OperationLogService
}

// New 构造路由组装器。
func New(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	supplierHandler *handler.SupplierHandler,
	inventoryHandler *handler.InventoryHandler,
	purchaseHandler *handler.PurchaseHandler,
	statsHandler *handler.StatsHandler,
	logHandler *handler.LogHandler,
	rateLimiter *middleware.RateLimiter,
	operationLogService *service.OperationLogService,
) *Router {
	return &Router{
		cfg:                 cfg,
		authHandler:         authHandler,
		userHandler:         userHandler,
		supplierHandler:     supplierHandler,
		inventoryHandler:    inventoryHandler,
		purchaseHandler:     purchaseHandler,
		statsHandler:        statsHandler,
		logHandler:          logHandler,
		rateLimiter:         rateLimiter,
		operationLogService: operationLogService,
	}
}

// Setup 构建并返回 gin.Engine。
func (r *Router) Setup() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.ErrorHandler())
	engine.Use(middleware.Audit())
	if r.rateLimiter != nil {
		engine.Use(r.rateLimiter.Handler())
	}

	authCfg := middleware.AuthConfig{
		JWTSecret:    r.cfg.Auth.JWTSecret,
		APIKeySecret: r.cfg.Auth.APIKeySecret,
	}
	authMW := middleware.Auth(authCfg)

	api := engine.Group("/api/v1")
	api.Use(middleware.OperationLogger(r.operationLogService))

	r.registerAuthRoutes(api, authMW)
	r.registerUserRoutes(api, authMW)
	r.registerSupplierRoutes(api, authMW)
	r.registerInventoryRoutes(api, authMW)
	r.registerPurchaseRoutes(api, authMW)
	r.registerStatsRoutes(api, authMW)
	r.registerLogRoutes(api, authMW)

	r.registerPublicRoutes(engine)
	return engine
}

// registerPublicRoutes 注册无需认证的路由（健康检查、Swagger、OpenAPI）。
func (r *Router) registerPublicRoutes(engine *gin.Engine) {
	engine.GET("/healthz", r.authHandler.Health)
	engine.GET("/openapi.json", r.serveOpenAPI)
	registerSwaggerRoutes(engine)
}
