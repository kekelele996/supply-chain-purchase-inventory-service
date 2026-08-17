// 餐饮供应链 API 服务入口：加载配置、装配依赖并启动 HTTP 服务。
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/supplychain/supplychain-api/internal/config"
	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/database"
	"github.com/supplychain/supplychain-api/internal/handler"
	"github.com/supplychain/supplychain-api/internal/middleware"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/router"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

func main() {
	if err := run(); err != nil {
		util.Logger.Error("server fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx := context.Background()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		return err
	}
	util.Logger.Info(constants.LogTplDBConnected, "host", cfg.Database.Host, "db", cfg.Database.Name)
	if err := database.Migrate(db); err != nil {
		return err
	}
	util.Logger.Info(constants.LogTplDBMigrateDone)
	if err := database.Seed(ctx, db); err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		util.Logger.Warn("redis ping failed, rate limit degraded", "err", err)
	}
	util.Logger.Info(constants.LogTplRedisConnected, "addr", cfg.Redis.Addr)

	// repository 层
	userRepo := repository.NewUserRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	orderRepo := repository.NewPurchaseRepository(db)
	logRepo := repository.NewOperationLogRepository(db)

	// service 层
	tokenStore := service.NewRedisTokenStore(redisClient)
	authService := service.NewAuthService(userRepo, tokenStore, service.AuthServiceConfig{
		JWTSecret:       cfg.Auth.JWTSecret,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})
	userService := service.NewUserService(userRepo)
	supplierService := service.NewSupplierService(supplierRepo, inventoryRepo, orderRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, supplierRepo)
	purchaseService := service.NewPurchaseService(db, orderRepo, supplierRepo, inventoryRepo)
	statsService := service.NewStatsService(db)
	operationLogService := service.NewOperationLogService(logRepo)

	// handler 层
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	purchaseHandler := handler.NewPurchaseHandler(purchaseService)
	statsHandler := handler.NewStatsHandler(statsService)
	logHandler := handler.NewLogHandler(operationLogService)

	rateLimiter := middleware.NewRateLimiter(redisClient, cfg.Limiter.PerMinute)

	r := router.New(cfg, authHandler, userHandler, supplierHandler, inventoryHandler, purchaseHandler, statsHandler, logHandler, rateLimiter, operationLogService)
	engine := r.Setup()

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		util.Logger.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			util.Logger.Error("server shutdown error", "err", err)
		}
		close(done)
	}()

	util.Logger.Info(constants.LogTplServerStart, "port", cfg.Server.Port, "mode", cfg.Server.Mode)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	<-done
	util.Logger.Info(constants.LogTplServerStopped, "err", nil)
	return nil
}
