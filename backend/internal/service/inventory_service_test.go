package service

import (
	"context"
	"testing"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

func TestInventoryService_Create(t *testing.T) {
	supRepo := newFakeSupplierRepo()
	invRepo := newFakeInventoryRepo()
	svc := NewInventoryService(invRepo, supRepo)
	ctx := context.Background()
	sup := &model.Supplier{
		Name: "供应商1", ContactPerson: "张三", Phone: "13800000000",
		Address: "addr", Categories: model.StringList{"蔬菜"}, Status: constants.SupplierActive,
	}
	if err := supRepo.Create(ctx, sup); err != nil {
		t.Fatalf("create supplier: %v", err)
	}
	item, err := svc.Create(ctx, CreateInventoryParams{
		Name: "番茄", Category: string(constants.CategoryVegetable), SupplierID: sup.ID,
		BatchNo: "B-1", Quantity: 50, Unit: "kg", MinThreshold: 10,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
		StorageLocation: "A-1",
	})
	if err != nil {
		t.Fatalf("create inventory: %v", err)
	}
	if item.Status != constants.InventoryNormal {
		t.Fatalf("expected normal status, got %v", item.Status)
	}
	// 重复批次号
	if _, err := svc.Create(ctx, CreateInventoryParams{
		Name: "番茄2", Category: string(constants.CategoryVegetable), SupplierID: sup.ID,
		BatchNo: "B-1", Quantity: 50, Unit: "kg", MinThreshold: 10,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
		StorageLocation: "A-1",
	}); err == nil {
		t.Fatalf("expected duplicate batch error")
	}
}

func TestInventoryService_AdjustQuantity(t *testing.T) {
	supRepo := newFakeSupplierRepo()
	invRepo := newFakeInventoryRepo()
	svc := NewInventoryService(invRepo, supRepo)
	ctx := context.Background()
	sup := &model.Supplier{
		Name: "供应商2", ContactPerson: "李四", Phone: "13900000000",
		Address: "addr", Categories: model.StringList{"肉类"}, Status: constants.SupplierActive,
	}
	if err := supRepo.Create(ctx, sup); err != nil {
		t.Fatalf("create supplier: %v", err)
	}
	item, err := svc.Create(ctx, CreateInventoryParams{
		Name: "猪肉", Category: string(constants.CategoryMeat), SupplierID: sup.ID,
		BatchNo: "B-2", Quantity: 20, Unit: "kg", MinThreshold: 5,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
		StorageLocation: "B-1",
	})
	if err != nil {
		t.Fatalf("create inventory: %v", err)
	}
	updated, err := svc.AdjustQuantity(ctx, item.ID, -5)
	if err != nil {
		t.Fatalf("adjust: %v", err)
	}
	if updated.Quantity != 15 {
		t.Fatalf("expected 15, got %v", updated.Quantity)
	}
	if _, err := svc.AdjustQuantity(ctx, item.ID, -100); err == nil {
		t.Fatalf("expected insufficient error")
	} else if ae := util.AsAppError(err); ae == nil || ae.Code != constants.CodeInventoryInsufficient {
		t.Fatalf("expected insufficient error, got %v", err)
	}
}

func TestInventoryService_Alerts(t *testing.T) {
	invRepo := newFakeInventoryRepo()
	svc := NewInventoryService(invRepo, newFakeSupplierRepo())
	ctx := context.Background()
	_ = invRepo.Create(ctx, &model.InventoryItem{
		Name: "低库存", Category: constants.CategoryVegetable, SupplierID: 1,
		BatchNo: "B-3", Quantity: 2, Unit: "kg", MinThreshold: 10,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour), StorageLocation: "A-1",
		Status: constants.InventoryLow,
	})
	_ = invRepo.Create(ctx, &model.InventoryItem{
		Name: "正常", Category: constants.CategoryVegetable, SupplierID: 1,
		BatchNo: "B-4", Quantity: 100, Unit: "kg", MinThreshold: 10,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour), StorageLocation: "A-1",
		Status: constants.InventoryNormal,
	})
	list, err := svc.Alerts(ctx)
	if err != nil {
		t.Fatalf("alerts: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(list))
	}
}

var _ repository.InventoryRepository = (*fakeInventoryRepo)(nil)
