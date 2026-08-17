package service

import (
	"context"
	"testing"

)

func TestGet_ItemsHaveInventory(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := svc.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].InventoryItem == nil {
		t.Fatalf("items should have inventory item loaded: %+v", got.Items)
	}
}

func TestInventoryList_HasSupplier(t *testing.T) {
	svc, _, _, _ := newPurchaseEnv(t)
	ctx := context.Background()
	invSvc := NewInventoryService(svc.inventory, svc.suppliers)
	list, _, err := invSvc.List(ctx, 1, 10, "", "", 0, "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].Supplier == nil {
		t.Fatalf("inventory should have supplier loaded: %+v", list)
	}
}

func TestPurchaseList_HasCreator(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	if _, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID}); err != nil {
		t.Fatalf("create: %v", err)
	}
	list, _, err := svc.List(ctx, 1, 10, "", 0, "", "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].Creator == nil {
		t.Fatalf("orders should have creator loaded: %+v", list)
	}
}
