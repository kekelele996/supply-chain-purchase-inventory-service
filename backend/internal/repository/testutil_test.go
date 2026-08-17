package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/supplychain/supplychain-api/internal/model"
)

// newTestDB 打开内存 SQLite 并迁移全部模型。
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(
		&model.User{},
		&model.Supplier{},
		&model.InventoryItem{},
		&model.PurchaseOrder{},
		&model.PurchaseOrderItem{},
		&model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
