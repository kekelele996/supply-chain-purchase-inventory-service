package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
	"github.com/supplychain/supplychain-api/pkg/strutil"
)

// PurchaseService 采购单服务，承担采购单状态机与库存联动事务。
type PurchaseService struct {
	db        *gorm.DB
	orders    repository.PurchaseRepository
	suppliers repository.SupplierRepository
	inventory repository.InventoryRepository
}

// NewPurchaseService 构造采购单服务（复用 supplier/inventory 仓储完成校验与库存联动）。
func NewPurchaseService(db *gorm.DB, orders repository.PurchaseRepository, suppliers repository.SupplierRepository, inventory repository.InventoryRepository) *PurchaseService {
	return &PurchaseService{db: db, orders: orders, suppliers: suppliers, inventory: inventory}
}

// Create 创建采购单（含明细，事务内生成单号并校验供应商状态）。
func (s *PurchaseService) Create(ctx context.Context, req CreatePurchaseParams) (*model.PurchaseOrder, error) {
	sup, err := s.suppliers.FindByID(ctx, req.SupplierID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", req.SupplierID), err)
		}
		return nil, fmt.Errorf("create purchase order: %w", err)
	}
	if sup.Status != constants.SupplierActive {
		return nil, util.NewAppError(constants.CodeOrderSupplierSuspended, http.StatusBadRequest,
			fmt.Sprintf(constants.ErrTextSupplierState, sup.ID, sup.Status), nil)
	}
	items, total, err := s.buildItems(ctx, req.Items)
	if err != nil {
		return nil, err
	}
	order := &model.PurchaseOrder{
		SupplierID:  req.SupplierID,
		Status:      constants.OrderDraft,
		TotalAmount: total,
		CreatorID:   req.CreatorID,
		Notes:       req.Notes,
	}
	var created *model.PurchaseOrder
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		seq, cerr := s.orders.CountByPeriod(ctx, tx, startOfDay(time.Now()), startOfDay(time.Now().Add(24*time.Hour)))
		if cerr != nil {
			return fmt.Errorf("create purchase order: %w", cerr)
		}
		order.OrderNo = strutil.OrderNo(time.Now(), int(seq)+1)
		o, cerr := s.createInTx(ctx, tx, order, items)
		if cerr != nil {
			return cerr
		}
		created = o
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// Update 更新采购单（仅 draft/rejected 状态）。
func (s *PurchaseService) Update(ctx context.Context, id uint, req UpdatePurchaseParams, operatorID uint, operatorRole constants.UserRole) (*model.PurchaseOrder, error) {
	var updated *model.PurchaseOrder
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.ForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
			}
			return fmt.Errorf("update purchase order: %w", err)
		}
		if order.Status != constants.OrderDraft && order.Status != constants.OrderRejected {
			return util.NewAppError(constants.CodeOrderStateInvalid, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInvalidState, "采购单", id, order.Status, "更新"), nil)
		}
		if order.CreatorID != operatorID && operatorRole != constants.RoleManager && operatorRole != constants.RoleAdmin {
			return util.NewAppError(constants.CodeForbidden, http.StatusForbidden,
				fmt.Sprintf(constants.ErrTextForbidden, operatorRole, "更新采购单"), nil)
		}
		items, total, err := s.buildItemsWithTx(ctx, tx, req.Items)
		if err != nil {
			return err
		}
		order.SupplierID = req.SupplierID
		order.Notes = req.Notes
		order.TotalAmount = total
		order.Status = constants.OrderDraft
		if err := s.orders.UpdateWithItems(ctx, tx, order, items); err != nil {
			return fmt.Errorf("update purchase order: %w", err)
		}
		updated = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete 删除采购单（仅 draft 状态）。
func (s *PurchaseService) Delete(ctx context.Context, id uint, operatorID uint, operatorRole constants.UserRole) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.ForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
			}
			return fmt.Errorf("delete purchase order: %w", err)
		}
		if order.Status != constants.OrderDraft {
			return util.NewAppError(constants.CodeOrderStateInvalid, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInvalidState, "采购单", id, order.Status, "删除"), nil)
		}
		if order.CreatorID != operatorID && operatorRole != constants.RoleAdmin && operatorRole != constants.RoleManager {
			return util.NewAppError(constants.CodeForbidden, http.StatusForbidden,
				fmt.Sprintf(constants.ErrTextForbidden, operatorRole, "删除采购单"), nil)
		}
		return s.orders.Delete(ctx, tx, id)
	})
	if err != nil {
		return err
	}
	return nil
}

// Submit 提交审批：draft/rejected -> pending_approval。
func (s *PurchaseService) Submit(ctx context.Context, id uint) (*model.PurchaseOrder, error) {
	var updated *model.PurchaseOrder
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.ForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
			}
			return fmt.Errorf("submit purchase order: %w", err)
		}
		if order.Status != constants.OrderDraft && order.Status != constants.OrderRejected {
			return util.NewAppError(constants.CodeOrderStateInvalid, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInvalidState, "采购单", id, order.Status, "提交审批"), nil)
		}
		order.Status = constants.OrderPendingApproval
		if err := s.orders.UpdateStatus(ctx, tx, order); err != nil {
			return fmt.Errorf("submit purchase order: %w", err)
		}
		updated = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Approve 审批通过：pending_approval -> approved（manager+ 且不能审批自己的单）。
func (s *PurchaseService) Approve(ctx context.Context, id uint, approverID uint, approverRole constants.UserRole) (*model.PurchaseOrder, error) {
	if approverRole != constants.RoleManager && approverRole != constants.RoleAdmin {
		return nil, util.NewAppError(constants.CodeForbidden, http.StatusForbidden,
			fmt.Sprintf(constants.ErrTextForbidden, approverRole, "审批采购单"), nil)
	}
	var updated *model.PurchaseOrder
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.ForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
			}
			return fmt.Errorf("approve purchase order: %w", err)
		}
		if order.Status != constants.OrderPendingApproval {
			return util.NewAppError(constants.CodeOrderStateInvalid, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInvalidState, "采购单", id, order.Status, "审批通过"), nil)
		}
		if order.CreatorID == approverID {
			return util.NewAppError(constants.CodeOrderCannotSelfApprove, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextSelfApprove, approverRole, id), nil)
		}
		now := time.Now()
		order.Status = constants.OrderApproved
		order.ApproverID = &approverID
		order.ApprovedAt = &now
		if err := s.orders.UpdateStatus(ctx, tx, order); err != nil {
			return fmt.Errorf("approve purchase order: %w", err)
		}
		updated = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Reject 审批拒绝：pending_approval -> rejected。
func (s *PurchaseService) Reject(ctx context.Context, id uint, approverID uint, approverRole constants.UserRole) (*model.PurchaseOrder, error) {
	if approverRole != constants.RoleManager && approverRole != constants.RoleAdmin {
		return nil, util.NewAppError(constants.CodeForbidden, http.StatusForbidden,
			fmt.Sprintf(constants.ErrTextForbidden, approverRole, "拒绝采购单"), nil)
	}
	var updated *model.PurchaseOrder
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.ForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
			}
			return fmt.Errorf("reject purchase order: %w", err)
		}
		if order.Status != constants.OrderPendingApproval {
			return util.NewAppError(constants.CodeOrderStateInvalid, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInvalidState, "采购单", id, order.Status, "审批拒绝"), nil)
		}
		now := time.Now()
		order.Status = constants.OrderRejected
		order.ApproverID = &approverID
		order.ApprovedAt = &now
		if err := s.orders.UpdateStatus(ctx, tx, order); err != nil {
			return fmt.Errorf("reject purchase order: %w", err)
		}
		updated = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Complete 标记完成：approved -> completed，事务内用行锁更新库存余量。
func (s *PurchaseService) Complete(ctx context.Context, id uint) (*model.PurchaseOrder, int, error) {
	var updated *model.PurchaseOrder
	var updatedItems int
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.ForUpdate(ctx, tx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
			}
			return fmt.Errorf("complete purchase order: %w", err)
		}
		if order.Status != constants.OrderApproved {
			return util.NewAppError(constants.CodeOrderStateInvalid, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInvalidState, "采购单", id, order.Status, "标记完成"), nil)
		}
		var items []model.PurchaseOrderItem
		if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
			return fmt.Errorf("complete purchase order: load items: %w", err)
		}
		for _, it := range items {
			inv, ierr := s.inventory.ForUpdate(ctx, tx, it.InventoryItemID)
			if ierr != nil {
				if errors.Is(ierr, repository.ErrNotFound) {
					return util.NewAppError(constants.CodeInventoryNotFound, http.StatusNotFound,
						fmt.Sprintf(constants.ErrTextNotFound, "库存", it.InventoryItemID), ierr)
				}
				return fmt.Errorf("complete purchase order: %w", ierr)
			}
			inv.Quantity += it.Quantity
			inv.Status = computeInventoryStatus(inv.Quantity, inv.MinThreshold, inv.ExpiryDate)
			if err := tx.Save(inv).Error; err != nil {
				return fmt.Errorf("complete purchase order: update inventory: %w", err)
			}
			updatedItems++
		}
		now := time.Now()
		order.Status = constants.OrderCompleted
		order.CompletedAt = &now
		if err := s.orders.UpdateStatus(ctx, tx, order); err != nil {
			return fmt.Errorf("complete purchase order: %w", err)
		}
		updated = order
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return updated, updatedItems, nil
}

// Get 获取采购单详情（含明细）。
func (s *PurchaseService) Get(ctx context.Context, id uint) (*model.PurchaseOrder, error) {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeOrderNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "采购单", id), err)
		}
		return nil, fmt.Errorf("get purchase order: %w", err)
	}
	return order, nil
}

// List 分页查询采购单。
func (s *PurchaseService) List(ctx context.Context, page, pageSize int, status string, supplierID uint, startDate, endDate string) ([]model.PurchaseOrder, int64, error) {
	list, total, err := s.orders.List(ctx, page, pageSize, status, supplierID, startDate, endDate)
	if err != nil {
		return nil, 0, fmt.Errorf("list purchase orders: %w", err)
	}
	return list, total, nil
}

// CreatePurchaseParams 创建采购单参数。
type CreatePurchaseParams struct {
	SupplierID uint
	Notes      string
	Items      []PurchaseItemParams
	CreatorID  uint
}

// UpdatePurchaseParams 更新采购单参数。
type UpdatePurchaseParams struct {
	SupplierID uint
	Notes      string
	Items      []PurchaseItemParams
}

// PurchaseItemParams 采购明细参数。
type PurchaseItemParams struct {
	InventoryItemID uint
	Quantity        float64
	UnitPrice       float64
}

func (s *PurchaseService) buildItems(ctx context.Context, items []PurchaseItemParams) ([]model.PurchaseOrderItem, float64, error) {
	return s.buildItemsWithTx(ctx, nil, items)
}

func (s *PurchaseService) buildItemsWithTx(ctx context.Context, tx *gorm.DB, items []PurchaseItemParams) ([]model.PurchaseOrderItem, float64, error) {
	out := make([]model.PurchaseOrderItem, 0, len(items))
	var total float64
	for _, it := range items {
		inv, err := s.loadInventory(ctx, tx, it.InventoryItemID)
		if err != nil {
			return nil, 0, err
		}
		if inv.Status == constants.InventoryExpired {
			return nil, 0, util.NewAppError(constants.CodeInventoryExpired, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextNotFound, "库存", it.InventoryItemID), nil)
		}
		subtotal := round2(it.Quantity * it.UnitPrice)
		total += subtotal
		out = append(out, model.PurchaseOrderItem{
			InventoryItemID: it.InventoryItemID,
			Quantity:        it.Quantity,
			UnitPrice:       it.UnitPrice,
			Subtotal:        subtotal,
		})
	}
	return out, round2(total), nil
}

func (s *PurchaseService) loadInventory(ctx context.Context, tx *gorm.DB, id uint) (*model.InventoryItem, error) {
	if tx != nil {
		return s.inventory.ForUpdate(ctx, tx, id)
	}
	return s.inventory.FindByID(ctx, id)
}

func (s *PurchaseService) createInTx(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder, items []model.PurchaseOrderItem) (*model.PurchaseOrder, error) {
	if err := tx.Create(order).Error; err != nil {
		if isDuplicateErrMsg(err) {
			return nil, util.NewAppError(constants.CodeConflict, http.StatusConflict, "采购单号冲突，请重试", err)
		}
		return nil, fmt.Errorf("create purchase order: %w", err)
	}
	for i := range items {
		items[i].OrderID = order.ID
	}
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			return nil, fmt.Errorf("create purchase order items: %w", err)
		}
	}
	return order, nil
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
