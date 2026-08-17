package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

// InventoryRepository 库存仓储接口。
type InventoryRepository interface {
	Create(ctx context.Context, item *model.InventoryItem) error
	Update(ctx context.Context, item *model.InventoryItem) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.InventoryItem, error)
	FindByBatchNo(ctx context.Context, batchNo string) (*model.InventoryItem, error)
	List(ctx context.Context, page, pageSize int, status, category string, supplierID uint, name string) ([]model.InventoryItem, int64, error)
	ListBySupplier(ctx context.Context, supplierID uint, page, pageSize int) ([]model.InventoryItem, int64, error)
	ListAlerts(ctx context.Context) ([]model.InventoryItem, error)
	AdjustQuantity(ctx context.Context, id uint, delta float64, newStatus constants.InventoryStatus) (*model.InventoryItem, error)
	// ForUpdate 在事务内锁行读取。
	ForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.InventoryItem, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 构造库存仓储。
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(ctx context.Context, item *model.InventoryItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create inventory: %w", ErrDuplicateKey)
		}
		return fmt.Errorf("create inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepository) Update(ctx context.Context, item *model.InventoryItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("update inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&model.InventoryItem{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete inventory id=%d: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete inventory id=%d: %w", id, ErrNotFound)
	}
	return nil
}

func (r *inventoryRepository) FindByID(ctx context.Context, id uint) (*model.InventoryItem, error) {
	var item model.InventoryItem
	if err := r.db.WithContext(ctx).Preload("Supplier").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find inventory id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find inventory id=%d: %w", id, err)
	}
	return &item, nil
}

func (r *inventoryRepository) FindByBatchNo(ctx context.Context, batchNo string) (*model.InventoryItem, error) {
	var item model.InventoryItem
	if err := r.db.WithContext(ctx).Where("batch_no = ?", batchNo).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find inventory batch %q: %w", batchNo, ErrNotFound)
		}
		return nil, fmt.Errorf("find inventory batch %q: %w", batchNo, err)
	}
	return &item, nil
}

func (r *inventoryRepository) List(ctx context.Context, page, pageSize int, status, category string, supplierID uint, name string) ([]model.InventoryItem, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.InventoryItem{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if supplierID > 0 {
		q = q.Where("supplier_id = ?", supplierID)
	}
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count inventory: %w", err)
	}
	var list []model.InventoryItem
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list inventory: %w", err)
	}
	return list, total, nil
}

func (r *inventoryRepository) ListBySupplier(ctx context.Context, supplierID uint, page, pageSize int) ([]model.InventoryItem, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.InventoryItem{}).Where("supplier_id = ?", supplierID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count inventory by supplier: %w", err)
	}
	var list []model.InventoryItem
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list inventory by supplier: %w", err)
	}
	return list, total, nil
}

func (r *inventoryRepository) ListAlerts(ctx context.Context) ([]model.InventoryItem, error) {
	var list []model.InventoryItem
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("quantity < min_threshold OR expiry_date < ?", now.Format("2006-01-02")).
		Order("expiry_date ASC").
		Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("list inventory alerts: %w", err)
	}
	return list, nil
}

func (r *inventoryRepository) AdjustQuantity(ctx context.Context, id uint, delta float64, newStatus constants.InventoryStatus) (*model.InventoryItem, error) {
	var item model.InventoryItem
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clauseLocking).First(&item, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("adjust quantity inventory id=%d: %w", id, ErrNotFound)
			}
			return fmt.Errorf("adjust quantity inventory id=%d: %w", id, err)
		}
		if item.Quantity+delta < 0 {
			return fmt.Errorf("adjust quantity inventory id=%d: %w", id, ErrConflict)
		}
		item.Quantity += delta
		item.Status = newStatus
		if err := tx.Save(&item).Error; err != nil {
			return fmt.Errorf("adjust quantity inventory id=%d: %w", id, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *inventoryRepository) ForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.InventoryItem, error) {
	var item model.InventoryItem
	if err := tx.WithContext(ctx).Clauses(clauseLocking).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("for update inventory id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("for update inventory id=%d: %w", id, err)
	}
	return &item, nil
}
