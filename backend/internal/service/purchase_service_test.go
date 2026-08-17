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
	"github.com/supplychain/supplychain-api/internal/util"
)

func newPurchaseEnv(t *testing.T) (*PurchaseService, *model.Supplier, *model.InventoryItem, *model.User) {
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
		&model.User{}, &model.Supplier{}, &model.InventoryItem{},
		&model.PurchaseOrder{}, &model.PurchaseOrderItem{}, &model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	ctx := context.Background()
	userRepo := repository.NewUserRepository(db)
	supRepo := repository.NewSupplierRepository(db)
	invRepo := repository.NewInventoryRepository(db)
	orderRepo := repository.NewPurchaseRepository(db)

	user := &model.User{Username: "op1", PasswordHash: "h", Role: constants.RoleOperator}
	manager := &model.User{Username: "mgr1", PasswordHash: "h", Role: constants.RoleManager}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := userRepo.Create(ctx, manager); err != nil {
		t.Fatalf("create manager: %v", err)
	}
	sup := &model.Supplier{
		Name: "供应商", ContactPerson: "张三", Phone: "13800000000", Address: "addr",
		Categories: model.StringList{"蔬菜"}, Status: constants.SupplierActive,
	}
	if err := supRepo.Create(ctx, sup); err != nil {
		t.Fatalf("create supplier: %v", err)
	}
	item := &model.InventoryItem{
		Name: "番茄", Category: constants.CategoryVegetable, SupplierID: sup.ID,
		BatchNo: "B-100", Quantity: 100, Unit: "kg", MinThreshold: 10,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour), StorageLocation: "A-1",
		Status: constants.InventoryNormal,
	}
	if err := invRepo.Create(ctx, item); err != nil {
		t.Fatalf("create inventory: %v", err)
	}
	svc := NewPurchaseService(db, orderRepo, supRepo, invRepo)
	_ = manager
	return svc, sup, item, user
}

func TestPurchaseService_StateMachine(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()

	order, err := svc.Create(ctx, CreatePurchaseParams{
		SupplierID: sup.ID,
		Notes:      "测试采购",
		Items: []PurchaseItemParams{
			{InventoryItemID: item.ID, Quantity: 10, UnitPrice: 5.5},
		},
		CreatorID: user.ID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if order.OrderNo == "" || order.TotalAmount != 55 {
		t.Fatalf("unexpected order: no=%s total=%v", order.OrderNo, order.TotalAmount)
	}
	if order.Status != constants.OrderDraft {
		t.Fatalf("expected draft, got %v", order.Status)
	}

	// submit -> pending_approval
	submitted, err := svc.Submit(ctx, order.ID)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if submitted.Status != constants.OrderPendingApproval {
		t.Fatalf("expected pending_approval, got %v", submitted.Status)
	}

	// 普通操作员不能审批
	operatorRole := constants.RoleOperator
	if _, err := svc.Approve(ctx, order.ID, 999, operatorRole); err == nil {
		t.Fatalf("expected forbidden for operator approve")
	}

	// 创建人不能审批自己的单
	adminID := user.ID
	if _, err := svc.Approve(ctx, order.ID, adminID, constants.RoleAdmin); err == nil {
		t.Fatalf("expected self-approve error")
	} else if ae := util.AsAppError(err); ae == nil || ae.Code != constants.CodeOrderCannotSelfApprove {
		t.Fatalf("expected self approve error, got %v", err)
	}

	// manager 审批通过
	approved, err := svc.Approve(ctx, order.ID, 999, constants.RoleManager)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if approved.Status != constants.OrderApproved || approved.ApproverID == nil {
		t.Fatalf("unexpected approved: %+v", approved)
	}

	// 完成 -> 库存增加
	completed, itemsUpdated, err := svc.Complete(ctx, order.ID)
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if completed.Status != constants.OrderCompleted || itemsUpdated != 1 {
		t.Fatalf("unexpected completed: %+v updated=%d", completed, itemsUpdated)
	}
	got, err := svc.Get(ctx, order.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].Subtotal != 55 {
		t.Fatalf("unexpected items: %+v", got.Items)
	}
}

func TestPurchaseService_CreateWithSuspendedSupplier(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	// 变更供应商为暂停
	sup.Status = constants.SupplierSuspended
	if err := svc.suppliers.Update(ctx, sup); err != nil {
		t.Fatalf("update supplier: %v", err)
	}
	_, err := svc.Create(ctx, CreatePurchaseParams{
		SupplierID: sup.ID,
		Items:      []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 1, UnitPrice: 10}},
		CreatorID:  user.ID,
	})
	if err == nil {
		t.Fatalf("expected error for suspended supplier")
	}
	if ae := util.AsAppError(err); ae == nil || ae.Code != constants.CodeOrderSupplierSuspended {
		t.Fatalf("expected supplier suspended error, got %v", err)
	}
}

func TestPurchaseService_RejectThenEditAndResubmit(t *testing.T) {
	svc, sup, item, user := newPurchaseEnv(t)
	ctx := context.Background()
	order, err := svc.Create(ctx, CreatePurchaseParams{
		SupplierID: sup.ID,
		Items:      []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 5, UnitPrice: 10}},
		CreatorID:  user.ID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(ctx, order.ID); err != nil {
		t.Fatalf("submit: %v", err)
	}
	rejected, err := svc.Reject(ctx, order.ID, 999, constants.RoleManager)
	if err != nil {
		t.Fatalf("reject: %v", err)
	}
	if rejected.Status != constants.OrderRejected {
		t.Fatalf("expected rejected, got %v", rejected.Status)
	}
	// 已拒绝可编辑并重置为草稿
	updated, err := svc.Update(ctx, order.ID, UpdatePurchaseParams{
		SupplierID: sup.ID,
		Items:      []PurchaseItemParams{{InventoryItemID: item.ID, Quantity: 8, UnitPrice: 12}},
	}, user.ID, constants.RoleOperator)
	if err != nil {
		t.Fatalf("update after reject: %v", err)
	}
	if updated.Status != constants.OrderDraft || updated.TotalAmount != 96 {
		t.Fatalf("unexpected updated: %+v", updated)
	}
	// 草稿状态下不能删除？草稿可删除
	if err := svc.Delete(ctx, order.ID, user.ID, constants.RoleOperator); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get(ctx, order.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}
