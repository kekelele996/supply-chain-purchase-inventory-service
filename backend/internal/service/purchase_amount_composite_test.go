package service

import (
	"context"
	"testing"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

func TestPurchaseCreate_AmountPrecision(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{
		SupplierID: sup.ID,
		Items:      []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1.25, UnitPrice: 2.516}},
		CreatorID:  user.ID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := svc.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.TotalAmount != 3.15 {
		t.Errorf("total = %v, want 3.15", got.TotalAmount)
	}
	if len(got.Items) != 1 || got.Items[0].Subtotal != 3.15 {
		t.Errorf("items = %+v, want subtotal 3.15", got.Items)
	}
}

func TestPurchaseList_NewestFirst(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	o1, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create1: %v", err)
	}
	o2, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 2, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create2: %v", err)
	}
	list, _, err := svc.List(ctx, 1, 10, "", 0, "", "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
	if list[0].ID != o2.ID {
		t.Errorf("first id = %d, want %d", list[0].ID, o2.ID)
	}
	_ = o1
}

func TestQuantify_Decimals(t *testing.T) {
	if got := util.Quantify(5.5, "kg"); got != "5.5 kg" {
		t.Errorf("Quantify = %q, want 5.5 kg", got)
	}
}

func TestIsValidRejectedStatus(t *testing.T) {
	if !constants.IsValidPurchaseOrderStatus(constants.OrderRejected) {
		t.Error("rejected should be a valid purchase order status")
	}
}
