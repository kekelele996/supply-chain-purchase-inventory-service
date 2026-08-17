package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

// InventoryService 库存服务。
type InventoryService struct {
	inventory repository.InventoryRepository
	suppliers repository.SupplierRepository
}

// NewInventoryService 构造库存服务（复用供应商仓储校验供应商存在性）。
func NewInventoryService(inventory repository.InventoryRepository, suppliers repository.SupplierRepository) *InventoryService {
	return &InventoryService{inventory: inventory, suppliers: suppliers}
}

// Create 新增库存（入库）。
func (s *InventoryService) Create(ctx context.Context, req CreateInventoryParams) (*model.InventoryItem, error) {
	if _, err := s.suppliers.FindByID(ctx, req.SupplierID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", req.SupplierID), err)
		}
		return nil, fmt.Errorf("create inventory: %w", err)
	}
	if _, err := s.inventory.FindByBatchNo(ctx, req.BatchNo); err == nil {
		return nil, util.NewAppError(constants.CodeInventoryBatchExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextBatchExists, req.BatchNo), nil)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("create inventory: %w", err)
	}
	expiry, err := time.ParseInLocation("2006-01-02", req.ExpiryDate, time.Local)
	if err != nil {
		return nil, util.NewAppError(constants.CodeValidationError, 422, fmt.Sprintf(constants.ErrTextInvalidParam, "expiry_date"), err)
	}
	item := &model.InventoryItem{
		Name:           req.Name,
		Category:       constants.InventoryCategory(req.Category),
		SupplierID:     req.SupplierID,
		BatchNo:        req.BatchNo,
		Quantity:       req.Quantity,
		Unit:           req.Unit,
		MinThreshold:   req.MinThreshold,
		ExpiryDate:     expiry,
		StorageLocation: req.StorageLocation,
		Status:         computeInventoryStatus(req.Quantity, req.MinThreshold, expiry),
	}
	if err := s.inventory.Create(ctx, item); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, util.NewAppError(constants.CodeInventoryBatchExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextBatchExists, req.BatchNo), err)
		}
		return nil, fmt.Errorf("create inventory: %w", err)
	}
	return item, nil
}

// Update 更新库存信息（余量不变，仅基本信息）。
func (s *InventoryService) Update(ctx context.Context, id uint, req UpdateInventoryParams) (*model.InventoryItem, error) {
	item, err := s.inventory.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeInventoryNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "库存", id), err)
		}
		return nil, fmt.Errorf("update inventory: %w", err)
	}
	if _, err := s.suppliers.FindByID(ctx, req.SupplierID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", req.SupplierID), err)
		}
		return nil, fmt.Errorf("update inventory: %w", err)
	}
	expiry, err := time.ParseInLocation("2006-01-02", req.ExpiryDate, time.Local)
	if err != nil {
		return nil, util.NewAppError(constants.CodeValidationError, 422, fmt.Sprintf(constants.ErrTextInvalidParam, "expiry_date"), err)
	}
	item.Name = req.Name
	item.Category = constants.InventoryCategory(req.Category)
	item.SupplierID = req.SupplierID
	item.Unit = req.Unit
	item.MinThreshold = req.MinThreshold
	item.ExpiryDate = expiry
	item.StorageLocation = req.StorageLocation
	item.Status = computeInventoryStatus(item.Quantity, item.MinThreshold, item.ExpiryDate)
	if err := s.inventory.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("update inventory: %w", err)
	}
	return item, nil
}

// Delete 删除库存记录（admin）。
func (s *InventoryService) Delete(ctx context.Context, id uint) error {
	if err := s.inventory.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeInventoryNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "库存", id), err)
		}
		return fmt.Errorf("delete inventory: %w", err)
	}
	return nil
}

// Get 获取库存详情。
func (s *InventoryService) Get(ctx context.Context, id uint) (*model.InventoryItem, error) {
	item, err := s.inventory.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeInventoryNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "库存", id), err)
		}
		return nil, fmt.Errorf("get inventory: %w", err)
	}
	return item, nil
}

// List 分页查询库存。
func (s *InventoryService) List(ctx context.Context, page, pageSize int, status, category string, supplierID uint, name string) ([]model.InventoryItem, int64, error) {
	list, total, err := s.inventory.List(ctx, page, pageSize, status, category, supplierID, name)
	if err != nil {
		return nil, 0, fmt.Errorf("list inventory: %w", err)
	}
	return list, total, nil
}

// Alerts 获取库存预警列表。
func (s *InventoryService) Alerts(ctx context.Context) ([]model.InventoryItem, error) {
	list, err := s.inventory.ListAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list inventory alerts: %w", err)
	}
	return list, nil
}

// AdjustQuantity 调整库存余量（delta 正数为入库，负数为出库）。
func (s *InventoryService) AdjustQuantity(ctx context.Context, id uint, delta float64) (*model.InventoryItem, error) {
	if delta == 0 {
		return nil, util.NewAppError(constants.CodeBadRequest, http.StatusBadRequest, fmt.Sprintf(constants.ErrTextInvalidParam, "delta"), nil)
	}
	item, err := s.inventory.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeInventoryNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "库存", id), err)
		}
		return nil, fmt.Errorf("adjust inventory quantity: %w", err)
	}
	if item.Quantity+delta < 0 {
		return nil, util.NewAppError(constants.CodeInventoryInsufficient, http.StatusBadRequest,
			fmt.Sprintf(constants.ErrTextInsufficient, id, item.Quantity, item.Unit), nil)
	}
	newQty := item.Quantity + delta
	newStatus := computeInventoryStatus(newQty, item.MinThreshold, item.ExpiryDate)
	updated, err := s.inventory.AdjustQuantity(ctx, id, delta, newStatus)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeInventoryNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "库存", id), err)
		}
		if errors.Is(err, repository.ErrConflict) {
			return nil, util.NewAppError(constants.CodeInventoryInsufficient, http.StatusBadRequest,
				fmt.Sprintf(constants.ErrTextInsufficient, id, item.Quantity, item.Unit), err)
		}
		return nil, fmt.Errorf("adjust inventory quantity: %w", err)
	}
	return updated, nil
}

// CreateInventoryParams 新增库存参数。
type CreateInventoryParams struct {
	Name           string
	Category       string
	SupplierID     uint
	BatchNo        string
	Quantity       float64
	Unit           string
	MinThreshold   float64
	ExpiryDate     string
	StorageLocation string
}

// UpdateInventoryParams 更新库存参数。
type UpdateInventoryParams struct {
	Name           string
	Category       string
	SupplierID     uint
	Unit           string
	MinThreshold   float64
	ExpiryDate     string
	StorageLocation string
}

// computeInventoryStatus 根据余量/阈值/保质期计算库存状态。
// 保质期到期日当天视为仍在有效期内（"还有一天才过期"不应判为过期），
// 余量等于预警阈值时仍属正常，仅当严格低于阈值才预警。
func computeInventoryStatus(quantity, minThreshold float64, expiry time.Time) constants.InventoryStatus {
	if time.Now().After(expiry) {
		return constants.InventoryExpired
	}
	if quantity < minThreshold {
		return constants.InventoryLow
	}
	return constants.InventoryNormal
}
