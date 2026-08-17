package service

import (
	"context"
	"testing"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

func TestApprove_SelfApproveBlocked(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(ctx, order.ID); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if _, err := svc.Approve(ctx, order.ID, user.ID, constants.RoleAdmin); err == nil {
		t.Fatal("expected self-approve error")
	}
}

func TestReject_DraftBlocked(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Reject(ctx, order.ID, 999, constants.RoleManager); err == nil {
		t.Fatal("expected state invalid for draft reject")
	}
}

func TestSubmit_RejectedResubmitAllowed(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(ctx, order.ID); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if _, err := svc.Reject(ctx, order.ID, 999, constants.RoleManager); err != nil {
		t.Fatalf("reject: %v", err)
	}
	if _, err := svc.Submit(ctx, order.ID); err != nil {
		t.Fatalf("resubmit rejected order: %v", err)
	}
}

func TestOrderStatusText_Approved(t *testing.T) {
	if got := util.OrderStatusText(constants.OrderApproved); got != "已批准" {
		t.Errorf("OrderStatusText(approved) = %q, want 已批准", got)
	}
}

func TestIsValidCompletedStatus(t *testing.T) {
	if !constants.IsValidPurchaseOrderStatus(constants.OrderCompleted) {
		t.Error("completed should be a valid purchase order status")
	}
}

func TestGet_HasApprover(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{SupplierID: sup.ID, Items: []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}}, CreatorID: user.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(ctx, order.ID); err != nil {
		t.Fatalf("submit: %v", err)
	}
	approved, err := svc.Approve(ctx, order.ID, 2, constants.RoleManager)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	got, err := svc.Get(ctx, approved.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Approver == nil {
		t.Fatal("approver should be loaded")
	}
}
