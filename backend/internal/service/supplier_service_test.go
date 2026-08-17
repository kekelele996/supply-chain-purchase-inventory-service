package service

import (
	"context"
	"testing"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/util"
)

func TestSupplierService_CreateAndDuplicate(t *testing.T) {
	svc := NewSupplierService(newFakeSupplierRepo(), newFakeInventoryRepo(), newFakePurchaseRepo())
	ctx := context.Background()
	sup, err := svc.Create(ctx, CreateSupplierParams{
		Name:          "测试供应商",
		ContactPerson: "张三",
		Phone:         "13800000000",
		Address:       "某市某路1号",
		Categories:    []string{"蔬菜", "肉类"},
		Rating:        4.5,
		Status:        string(constants.SupplierActive),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if sup.ID == 0 || sup.Status != constants.SupplierActive {
		t.Fatalf("unexpected supplier: %+v", sup)
	}
	if _, err := svc.Create(ctx, CreateSupplierParams{
		Name: "测试供应商", ContactPerson: "李四", Phone: "13900000000",
		Address: "某市某路2号", Categories: []string{"蔬菜"},
	}); err == nil {
		t.Fatalf("expected duplicate name error")
	} else if ae := util.AsAppError(err); ae == nil || ae.Code != constants.CodeSupplierNameExists {
		t.Fatalf("expected name exists error, got %v", err)
	}
}

func TestSupplierService_UpdateAndStatus(t *testing.T) {
	svc := NewSupplierService(newFakeSupplierRepo(), newFakeInventoryRepo(), newFakePurchaseRepo())
	ctx := context.Background()
	sup, err := svc.Create(ctx, CreateSupplierParams{
		Name: "A", ContactPerson: "张三", Phone: "13800000000",
		Address: "addr", Categories: []string{"蔬菜"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	updated, err := svc.Update(ctx, sup.ID, UpdateSupplierParams{
		Name: "A2", ContactPerson: "李四", Phone: "13900000000",
		Address: "addr2", Categories: []string{"海鲜"},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "A2" || len(updated.Categories) != 1 || updated.Categories[0] != "海鲜" {
		t.Fatalf("unexpected updated supplier: %+v", updated)
	}
	changed, err := svc.ChangeStatus(ctx, sup.ID, constants.SupplierSuspended)
	if err != nil {
		t.Fatalf("change status: %v", err)
	}
	if changed.Status != constants.SupplierSuspended {
		t.Fatalf("unexpected status: %v", changed.Status)
	}
	if err := svc.Delete(ctx, sup.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get(ctx, sup.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}

func TestSupplierService_ListInventoryAndOrders(t *testing.T) {
	supRepo := newFakeSupplierRepo()
	invRepo := newFakeInventoryRepo()
	orderRepo := newFakePurchaseRepo()
	svc := NewSupplierService(supRepo, invRepo, orderRepo)
	ctx := context.Background()
	sup, _ := svc.Create(ctx, CreateSupplierParams{
		Name: "B", ContactPerson: "王五", Phone: "13700000000",
		Address: "addr", Categories: []string{"肉类"},
	})
	if err := invRepo.Create(ctx, &model.InventoryItem{
		Name: "食材", Category: constants.CategoryMeat, SupplierID: sup.ID,
		BatchNo: "X1", Quantity: 10, Unit: "kg", MinThreshold: 5,
		ExpiryDate: time.Now().Add(10 * 24 * time.Hour), StorageLocation: "A-1",
		Status: constants.InventoryNormal,
	}); err != nil {
		t.Fatalf("create inventory: %v", err)
	}
	list, total, err := svc.ListInventory(ctx, sup.ID, 1, 10)
	if err != nil {
		t.Fatalf("list inventory: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("expected 1 inventory, got %d", total)
	}
	orders, total, err := svc.ListOrders(ctx, sup.ID, 1, 10)
	if err != nil {
		t.Fatalf("list orders: %v", err)
	}
	if total != 0 || len(orders) != 0 {
		t.Fatalf("expected 0 orders, got %d", total)
	}
}
