package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/model"
)

// SupplierRepository 供应商仓储接口。
type SupplierRepository interface {
	Create(ctx context.Context, s *model.Supplier) error
	Update(ctx context.Context, s *model.Supplier) error
	SoftDelete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Supplier, error)
	FindByName(ctx context.Context, name string) (*model.Supplier, error)
	List(ctx context.Context, page, pageSize int, status, category, name string) ([]model.Supplier, int64, error)
}

type supplierRepository struct {
	db *gorm.DB
}

// NewSupplierRepository 构造供应商仓储。
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, s *model.Supplier) error {
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create supplier: %w", ErrDuplicateKey)
		}
		return fmt.Errorf("create supplier: %w", err)
	}
	return nil
}

func (r *supplierRepository) Update(ctx context.Context, s *model.Supplier) error {
	if err := r.db.WithContext(ctx).Save(s).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("update supplier: %w", ErrDuplicateKey)
		}
		return fmt.Errorf("update supplier: %w", err)
	}
	return nil
}

func (r *supplierRepository) SoftDelete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Supplier{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete supplier id=%d: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete supplier id=%d: %w", id, ErrNotFound)
	}
	return nil
}

func (r *supplierRepository) FindByID(ctx context.Context, id uint) (*model.Supplier, error) {
	var s model.Supplier
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find supplier id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find supplier id=%d: %w", id, err)
	}
	return &s, nil
}

func (r *supplierRepository) FindByName(ctx context.Context, name string) (*model.Supplier, error) {
	var s model.Supplier
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find supplier %q: %w", name, ErrNotFound)
		}
		return nil, fmt.Errorf("find supplier %q: %w", name, err)
	}
	return &s, nil
}

func (r *supplierRepository) List(ctx context.Context, page, pageSize int, status, category, name string) ([]model.Supplier, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Supplier{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		// 使用 LIKE 匹配 JSON 数组中的字符串元素，兼容 MySQL 与 SQLite。
		q = q.Where("categories LIKE ?", "%"+category+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count suppliers: %w", err)
	}
	var list []model.Supplier
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list suppliers: %w", err)
	}
	return list, total, nil
}
