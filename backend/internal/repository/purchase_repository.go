package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/model"
)

// PurchaseRepository 采购单仓储接口。
type PurchaseRepository interface {
	Create(ctx context.Context, order *model.PurchaseOrder, items []model.PurchaseOrderItem) error
	UpdateWithItems(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder, items []model.PurchaseOrderItem) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	FindByID(ctx context.Context, id uint) (*model.PurchaseOrder, error)
	FindByOrderNo(ctx context.Context, orderNo string) (*model.PurchaseOrder, error)
	List(ctx context.Context, page, pageSize int, status string, supplierID uint, startDate, endDate string) ([]model.PurchaseOrder, int64, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder) error
	// ForUpdate 在事务内锁行读取采购单。
	ForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.PurchaseOrder, error)
	// CountByPeriod 统计某段时间内已创建的单量（用于生成单号序号）。
	CountByPeriod(ctx context.Context, tx *gorm.DB, start, end time.Time) (int64, error)
}

type purchaseRepository struct {
	db *gorm.DB
}

// NewPurchaseRepository 构造采购单仓储。
func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) Create(ctx context.Context, order *model.PurchaseOrder, items []model.PurchaseOrderItem) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			if isDuplicateErr(err) {
				return fmt.Errorf("create purchase order: %w", ErrDuplicateKey)
			}
			return fmt.Errorf("create purchase order: %w", err)
		}
		for i := range items {
			items[i].OrderID = order.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return fmt.Errorf("create purchase order items: %w", err)
			}
		}
		return nil
	})
	return err
}

func (r *purchaseRepository) UpdateWithItems(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder, items []model.PurchaseOrderItem) error {
	if err := tx.WithContext(ctx).Save(order).Error; err != nil {
		return fmt.Errorf("update purchase order: %w", err)
	}
	if err := tx.WithContext(ctx).Where("order_id = ?", order.ID).Delete(&model.PurchaseOrderItem{}).Error; err != nil {
		return fmt.Errorf("delete purchase order items: %w", err)
	}
	for i := range items {
		items[i].ID = 0
		items[i].OrderID = order.ID
	}
	if len(items) > 0 {
		if err := tx.WithContext(ctx).Create(&items).Error; err != nil {
			return fmt.Errorf("create purchase order items: %w", err)
		}
	}
	return nil
}

func (r *purchaseRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := tx.WithContext(ctx).Where("order_id = ?", id).Delete(&model.PurchaseOrderItem{}).Error; err != nil {
		return fmt.Errorf("delete purchase order items: %w", err)
	}
	res := tx.WithContext(ctx).Delete(&model.PurchaseOrder{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete purchase order: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete purchase order id=%d: %w", id, ErrNotFound)
	}
	return nil
}

func (r *purchaseRepository) FindByID(ctx context.Context, id uint) (*model.PurchaseOrder, error) {
	var order model.PurchaseOrder
	if err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Creator").
		Preload("Approver").
		Preload("Items.InventoryItem").
		First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find purchase order id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find purchase order id=%d: %w", id, err)
	}
	return &order, nil
}

func (r *purchaseRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.PurchaseOrder, error) {
	var order model.PurchaseOrder
	if err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find purchase order %q: %w", orderNo, ErrNotFound)
		}
		return nil, fmt.Errorf("find purchase order %q: %w", orderNo, err)
	}
	return &order, nil
}

func (r *purchaseRepository) List(ctx context.Context, page, pageSize int, status string, supplierID uint, startDate, endDate string) ([]model.PurchaseOrder, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.PurchaseOrder{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if supplierID > 0 {
		q = q.Where("supplier_id = ?", supplierID)
	}
	if startDate != "" {
		q = q.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		q = q.Where("created_at <= ?", endDate+" 23:59:59")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count purchase orders: %w", err)
	}
	var list []model.PurchaseOrder
	if err := q.Preload("Supplier").Preload("Creator").
		Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list purchase orders: %w", err)
	}
	return list, total, nil
}

func (r *purchaseRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder) error {
	if err := tx.WithContext(ctx).Model(order).Updates(map[string]interface{}{
		"status":       order.Status,
		"approver_id":  order.ApproverID,
		"approved_at":  order.ApprovedAt,
		"completed_at": order.CompletedAt,
	}).Error; err != nil {
		return fmt.Errorf("update purchase order status: %w", err)
	}
	return nil
}

func (r *purchaseRepository) ForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.PurchaseOrder, error) {
	var order model.PurchaseOrder
	if err := tx.WithContext(ctx).Clauses(clauseLocking).First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("for update purchase order id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("for update purchase order id=%d: %w", id, err)
	}
	return &order, nil
}

func (r *purchaseRepository) CountByPeriod(ctx context.Context, tx *gorm.DB, start, end time.Time) (int64, error) {
	var count int64
	if err := tx.WithContext(ctx).Model(&model.PurchaseOrder{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count purchase orders by period: %w", err)
	}
	return count, nil
}
