package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/model"
)

// OperationLogRepository 操作日志仓储接口。
type OperationLogRepository interface {
	Create(ctx context.Context, log *model.OperationLog) error
	FindByID(ctx context.Context, id uint) (*model.OperationLog, error)
	List(ctx context.Context, page, pageSize int, userID uint, action, targetType, startDate, endDate string) ([]model.OperationLog, int64, error)
}

type operationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository 构造操作日志仓储。
func NewOperationLogRepository(db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(ctx context.Context, log *model.OperationLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create operation log: %w", err)
	}
	return nil
}

func (r *operationLogRepository) FindByID(ctx context.Context, id uint) (*model.OperationLog, error) {
	var log model.OperationLog
	if err := r.db.WithContext(ctx).Preload("User").First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find operation log id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find operation log id=%d: %w", id, err)
	}
	return &log, nil
}

func (r *operationLogRepository) List(ctx context.Context, page, pageSize int, userID uint, action, targetType, startDate, endDate string) ([]model.OperationLog, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.OperationLog{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if action != "" {
		q = q.Where("action = ?", action)
	}
	if targetType != "" {
		q = q.Where("target_type = ?", targetType)
	}
	if startDate != "" {
		q = q.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		q = q.Where("created_at <= ?", endDate+" 23:59:59")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count operation logs: %w", err)
	}
	var list []model.OperationLog
	if err := q.Preload("User").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list operation logs: %w", err)
	}
	return list, total, nil
}
