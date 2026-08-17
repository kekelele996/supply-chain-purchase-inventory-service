package service

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
)

func newAlertStatsEnv(t *testing.T) (*gorm.DB, *PurchaseService, *InventoryService, *model.Supplier, *model.InventoryItem, *model.User) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&model.User{}, &model.Supplier{}, &model.InventoryItem{}, &model.PurchaseOrder{}, &model.PurchaseOrderItem{}, &model.OperationLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	ctx := context.Background()
	userRepo := repository.NewUserRepository(db)
	supRepo := repository.NewSupplierRepository(db)
	invRepo := repository.NewInventoryRepository(db)
	orderRepo := repository.NewPurchaseRepository(db)
	user := &model.User{Username: "op1", PasswordHash: "h", Role: constants.RoleOperator}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	sup := &model.Supplier{Name: "供应商", ContactPerson: "张三", Phone: "13800000000", Address: "addr", Categories: model.StringList{"蔬菜"}, Status: constants.SupplierActive}
	if err := supRepo.Create(ctx, sup); err != nil {
		t.Fatalf("create supplier: %v", err)
	}
	item := &model.InventoryItem{Name: "番茄", Category: constants.CategoryVegetable, SupplierID: sup.ID, BatchNo: "B-100", Quantity: 100, Unit: "kg", MinThreshold: 10, ExpiryDate: time.Now().Add(30 * 24 * time.Hour), StorageLocation: "A-1", Status: constants.InventoryNormal}
	if err := invRepo.Create(ctx, item); err != nil {
		t.Fatalf("create inventory: %v", err)
	}
	svc := NewPurchaseService(db, orderRepo, supRepo, invRepo)
	invSvc := NewInventoryService(invRepo, supRepo)
	return db, svc, invSvc, sup, item, user
}

func TestAlerts_IncludeExpired(t *testing.T) {
	_, _, invSvc, sup, _, _ := newAlertStatsEnv(t)
	ctx := context.Background()
	exp := &model.InventoryItem{Name: "过期菜", Category: constants.CategoryVegetable, SupplierID: sup.ID, BatchNo: "B-EXP", Quantity: 100, Unit: "kg", MinThreshold: 10, ExpiryDate: time.Now().Add(-24 * time.Hour), StorageLocation: "A-2", Status: constants.InventoryExpired}
	if err := invSvc.inventory.Create(ctx, exp); err != nil {
		t.Fatalf("create expired item: %v", err)
	}
	list, err := invSvc.Alerts(ctx)
	if err != nil {
		t.Fatalf("alerts: %v", err)
	}
	for _, it := range list {
		if it.ID == exp.ID {
			return
		}
	}
	t.Fatal("expired item should be included in alerts")
}

func TestStatsSummary_IncludeCompleted(t *testing.T) {
	db, svc, _, sup, item, user := newAlertStatsEnv(t)
	ctx := context.Background()
	o1, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 10, UnitPrice: 5}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create1: %v", err)
	}
	if _, err := svc.Submit(ctx, o1.ID); err != nil {
		t.Fatalf("submit1: %v", err)
	}
	if _, err := svc.Approve(ctx, o1.ID, 999, constants.RoleManager); err != nil {
		t.Fatalf("approve1: %v", err)
	}
	if _, _, err := svc.Complete(ctx, o1.ID); err != nil {
		t.Fatalf("complete1: %v", err)
	}
	o2, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 3, UnitPrice: 4}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create2: %v", err)
	}
	if _, err := svc.Submit(ctx, o2.ID); err != nil {
		t.Fatalf("submit2: %v", err)
	}
	if _, err := svc.Approve(ctx, o2.ID, 999, constants.RoleManager); err != nil {
		t.Fatalf("approve2: %v", err)
	}
	stats := NewStatsService(db)
	sum, err := stats.Summary(ctx)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	want := o1.TotalAmount + o2.TotalAmount
	if got, _ := sum["total_purchase_amount"].(float64); got != want {
		t.Errorf("total = %v, want %v", got, want)
	}
}
