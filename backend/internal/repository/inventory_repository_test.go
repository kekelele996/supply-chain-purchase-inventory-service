package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

func newInventory(supplierID uint, batch string, qty float64, expiry time.Time) *model.InventoryItem {
	return &model.InventoryItem{
		Name:            "测试食材",
		Category:        constants.CategoryVegetable,
		SupplierID:      supplierID,
		BatchNo:         batch,
		Quantity:        qty,
		Unit:            "kg",
		MinThreshold:    10,
		ExpiryDate:      expiry,
		StorageLocation: "A-1",
		Status:          constants.InventoryNormal,
	}
}

func TestInventoryRepository_CreateAndAdjust(t *testing.T) {
	db := newTestDB(t)
	repo := NewInventoryRepository(db)
	ctx := context.Background()
	item := newInventory(1, "B-001", 50, time.Now().Add(30*24*time.Hour))
	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Quantity != 50 {
		t.Fatalf("unexpected qty: %v", got.Quantity)
	}
	updated, err := repo.AdjustQuantity(ctx, item.ID, -10, constants.InventoryNormal)
	if err != nil {
		t.Fatalf("adjust: %v", err)
	}
	if updated.Quantity != 40 {
		t.Fatalf("expected 40 after adjust, got %v", updated.Quantity)
	}
}

func TestInventoryRepository_AdjustBelowZero(t *testing.T) {
	db := newTestDB(t)
	repo := NewInventoryRepository(db)
	ctx := context.Background()
	item := newInventory(1, "B-002", 5, time.Now().Add(30*24*time.Hour))
	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := repo.AdjustQuantity(ctx, item.ID, -10, constants.InventoryNormal); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestInventoryRepository_ListAlerts(t *testing.T) {
	db := newTestDB(t)
	repo := NewInventoryRepository(db)
	ctx := context.Background()
	low := newInventory(1, "B-003", 5, time.Now().Add(30*24*time.Hour))
	low.Status = constants.InventoryLow
	if err := repo.Create(ctx, low); err != nil {
		t.Fatalf("create low: %v", err)
	}
	expired := newInventory(1, "B-004", 50, time.Now().Add(-24*time.Hour))
	expired.Status = constants.InventoryExpired
	if err := repo.Create(ctx, expired); err != nil {
		t.Fatalf("create expired: %v", err)
	}
	normal := newInventory(1, "B-005", 50, time.Now().Add(30*24*time.Hour))
	if err := repo.Create(ctx, normal); err != nil {
		t.Fatalf("create normal: %v", err)
	}
	list, err := repo.ListAlerts(ctx)
	if err != nil {
		t.Fatalf("list alerts: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(list))
	}
}

func TestInventoryRepository_DuplicateBatch(t *testing.T) {
	db := newTestDB(t)
	repo := NewInventoryRepository(db)
	ctx := context.Background()
	a := newInventory(1, "B-DUP", 50, time.Now().Add(30*24*time.Hour))
	b := newInventory(1, "B-DUP", 50, time.Now().Add(30*24*time.Hour))
	if err := repo.Create(ctx, a); err != nil {
		t.Fatalf("create a: %v", err)
	}
	if err := repo.Create(ctx, b); !errors.Is(err, ErrDuplicateKey) {
		t.Fatalf("expected duplicate key, got %v", err)
	}
}
