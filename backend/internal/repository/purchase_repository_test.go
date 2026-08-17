package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

func TestPurchaseRepository_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := NewPurchaseRepository(db)
	ctx := context.Background()

	order := &model.PurchaseOrder{
		OrderNo:     "PO-20260817-0001",
		SupplierID:  1,
		Status:      constants.OrderDraft,
		TotalAmount: 120.5,
		CreatorID:   1,
	}
	items := []model.PurchaseOrderItem{
		{InventoryItemID: 1, Quantity: 10, UnitPrice: 5.0, Subtotal: 50.0},
		{InventoryItemID: 2, Quantity: 3, UnitPrice: 23.5, Subtotal: 70.5},
	}
	if err := repo.Create(ctx, order, items); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.OrderNo != "PO-20260817-0001" || len(got.Items) != 2 {
		t.Fatalf("unexpected order: %+v", got)
	}
	if got.TotalAmount != 120.5 {
		t.Fatalf("unexpected total: %v", got.TotalAmount)
	}
}

func TestPurchaseRepository_UpdateWithItemsAndDelete(t *testing.T) {
	db := newTestDB(t)
	repo := NewPurchaseRepository(db)
	ctx := context.Background()

	order := &model.PurchaseOrder{
		OrderNo:     "PO-20260817-0002",
		SupplierID:  1,
		Status:      constants.OrderDraft,
		TotalAmount: 10,
		CreatorID:   1,
	}
	items := []model.PurchaseOrderItem{{InventoryItemID: 1, Quantity: 1, UnitPrice: 10, Subtotal: 10}}
	if err := repo.Create(ctx, order, items); err != nil {
		t.Fatalf("create: %v", err)
	}
	newItems := []model.PurchaseOrderItem{{InventoryItemID: 2, Quantity: 2, UnitPrice: 20, Subtotal: 40}}
	order.TotalAmount = 40
	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.UpdateWithItems(ctx, tx, order, newItems)
	})
	if err != nil {
		t.Fatalf("update with items: %v", err)
	}
	got, err := repo.FindByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].Quantity != 2 {
		t.Fatalf("unexpected items after update: %+v", got.Items)
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		return repo.Delete(ctx, tx, order.ID)
	})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, order.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPurchaseRepository_CountByPeriod(t *testing.T) {
	db := newTestDB(t)
	repo := NewPurchaseRepository(db)
	ctx := context.Background()
	now := time.Now()
	for i := 0; i < 3; i++ {
		o := &model.PurchaseOrder{
			OrderNo:    "PO-X-" + string(rune('0'+i)),
			SupplierID: 1,
			Status:     constants.OrderDraft,
			CreatorID:  1,
		}
		if err := repo.Create(ctx, o, nil); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	err := db.Transaction(func(tx *gorm.DB) error {
		count, cerr := repo.CountByPeriod(ctx, tx, start, end)
		if cerr != nil {
			return cerr
		}
		if count != 3 {
			t.Fatalf("expected 3, got %d", count)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("count by period: %v", err)
	}
}
