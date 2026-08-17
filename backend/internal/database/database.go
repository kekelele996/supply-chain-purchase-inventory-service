// Package database 负责 MySQL 连接、自动迁移与种子数据。
package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/supplychain/supplychain-api/internal/config"
	"github.com/supplychain/supplychain-api/internal/model"
)

// Connect 建立 MySQL 连接，带重试等待（容器启动时数据库可能尚未就绪）。
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()
	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Warn),
		})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

// Migrate 执行 GORM AutoMigrate 创建/更新表结构。
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Supplier{},
		&model.InventoryItem{},
		&model.PurchaseOrder{},
		&model.PurchaseOrderItem{},
		&model.OperationLog{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

// Ping 校验数据库连通性。
func Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}
