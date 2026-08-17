package service

import (
	"context"
	"testing"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

func TestComputeStatus_AtThreshold(t *testing.T) {
	if got := computeInventoryStatus(10, 10, time.Now().Add(24*time.Hour)); got != constants.InventoryNormal {
		t.Errorf("at threshold = %s, want normal", got)
	}
}

func TestComputeStatus_TomorrowNotExpired(t *testing.T) {
	if got := computeInventoryStatus(100, 10, time.Now().Add(24*time.Hour)); got != constants.InventoryNormal {
		t.Errorf("tomorrow expiry = %s, want normal", got)
	}
}

func TestComplete_UpdatesInventoryStatus(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	item.Quantity = 5
	item.Status = constants.InventoryLow
	if err := svc.inventory.Update(ctx, item); err != nil {
		t.Fatalf("update item: %v", err)
	}
	order, err := svc.Create(ctx, CreatePurchaseParams{
		SupplierID: sup.ID,
		Items:      []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 10, UnitPrice: 5}},
		CreatorID:  user.ID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(ctx, order.ID); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if _, err := svc.Approve(ctx, order.ID, 999, constants.RoleManager); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if _, _, err := svc.Complete(ctx, order.ID); err != nil {
		t.Fatalf("complete: %v", err)
	}
	got, err := svc.inventory.FindByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("find item: %v", err)
	}
	if got.Quantity != 15 {
		t.Errorf("quantity = %v, want 15", got.Quantity)
	}
	if got.Status != constants.InventoryNormal {
		t.Errorf("status = %s, want normal", got.Status)
	}
}

func TestInventoryStatusText_Normal(t *testing.T) {
	if got := util.InventoryStatusText(constants.InventoryNormal); got != "正常" {
		t.Errorf("InventoryStatusText(normal) = %q, want 正常", got)
	}
}
