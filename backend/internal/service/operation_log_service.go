package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

// OperationLogService 操作日志服务。
type OperationLogService struct {
	logs repository.OperationLogRepository
}

// NewOperationLogService 构造操作日志服务。
func NewOperationLogService(logs repository.OperationLogRepository) *OperationLogService {
	return &OperationLogService{logs: logs}
}

// RecordOperation 记录一条操作日志（action/target_type/target_id/detail）。
func (s *OperationLogService) RecordOperation(ctx context.Context, userID uint, action constants.OperationAction, targetType string, targetID uint, detail map[string]interface{}) {
	detailJSON := ""
	if detail != nil {
		if b, err := json.Marshal(detail); err == nil {
			detailJSON = string(b)
		}
	}
	entry := &model.OperationLog{
		UserID:     userID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detailJSON,
		CreatedAt:  time.Now(),
	}
	if err := s.logs.Create(ctx, entry); err != nil {
		util.L(ctx).Error("record operation log failed", "err", err)
	}
}

// List 分页查询操作日志。
func (s *OperationLogService) List(ctx context.Context, page, pageSize int, userID uint, action, targetType, startDate, endDate string) ([]model.OperationLog, int64, error) {
	list, total, err := s.logs.List(ctx, page, pageSize, userID, action, targetType, startDate, endDate)
	if err != nil {
		return nil, 0, fmt.Errorf("list operation logs: %w", err)
	}
	return list, total, nil
}

// Get 获取日志详情。
func (s *OperationLogService) Get(ctx context.Context, id uint) (*model.OperationLog, error) {
	entry, err := s.logs.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeLogNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "日志", id), err)
		}
		return nil, fmt.Errorf("get operation log: %w", err)
	}
	return entry, nil
}
