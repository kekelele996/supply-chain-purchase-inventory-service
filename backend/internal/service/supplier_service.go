package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

// SupplierService 供应商服务。
type SupplierService struct {
	suppliers  repository.SupplierRepository
	inventory  repository.InventoryRepository
	orders     repository.PurchaseRepository
}

// NewSupplierService 构造供应商服务（复用 inventory/order 仓储查询关联数据）。
func NewSupplierService(suppliers repository.SupplierRepository, inventory repository.InventoryRepository, orders repository.PurchaseRepository) *SupplierService {
	return &SupplierService{suppliers: suppliers, inventory: inventory, orders: orders}
}

// Create 新增供应商。
func (s *SupplierService) Create(ctx context.Context, req CreateSupplierParams) (*model.Supplier, error) {
	if _, err := s.suppliers.FindByName(ctx, req.Name); err == nil {
		return nil, util.NewAppError(constants.CodeSupplierNameExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextNameExists, "供应商", req.Name), nil)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("create supplier: %w", err)
	}
	status := constants.SupplierActive
	if req.Status != "" {
		status = constants.SupplierStatus(req.Status)
	}
	rating := req.Rating
	if rating == 0 {
		rating = 3.0
	}
	sup := &model.Supplier{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Categories:    model.StringList(req.Categories),
		Rating:        rating,
		Status:        status,
	}
	if err := s.suppliers.Create(ctx, sup); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, util.NewAppError(constants.CodeSupplierNameExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextNameExists, "供应商", req.Name), err)
		}
		return nil, fmt.Errorf("create supplier: %w", err)
	}
	return sup, nil
}

// Update 更新供应商。
func (s *SupplierService) Update(ctx context.Context, id uint, req UpdateSupplierParams) (*model.Supplier, error) {
	sup, err := s.suppliers.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", id), err)
		}
		return nil, fmt.Errorf("update supplier: %w", err)
	}
	if existing, ferr := s.suppliers.FindByName(ctx, req.Name); ferr == nil && existing.ID != id {
		return nil, util.NewAppError(constants.CodeSupplierNameExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextNameExists, "供应商", req.Name), nil)
	} else if ferr != nil && !errors.Is(ferr, repository.ErrNotFound) {
		return nil, fmt.Errorf("update supplier: %w", ferr)
	}
	sup.Name = req.Name
	sup.ContactPerson = req.ContactPerson
	sup.Phone = req.Phone
	sup.Email = req.Email
	sup.Address = req.Address
	sup.Categories = model.StringList(req.Categories)
	if req.Rating > 0 {
		sup.Rating = req.Rating
	}
	if err := s.suppliers.Update(ctx, sup); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, util.NewAppError(constants.CodeSupplierNameExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextNameExists, "供应商", req.Name), err)
		}
		return nil, fmt.Errorf("update supplier: %w", err)
	}
	return sup, nil
}

// Delete 软删除供应商（admin）。
func (s *SupplierService) Delete(ctx context.Context, id uint) error {
	if err := s.suppliers.SoftDelete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", id), err)
		}
		return fmt.Errorf("delete supplier: %w", err)
	}
	return nil
}

// ChangeStatus 变更供应商状态（manager+）。
func (s *SupplierService) ChangeStatus(ctx context.Context, id uint, status constants.SupplierStatus) (*model.Supplier, error) {
	sup, err := s.suppliers.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", id), err)
		}
		return nil, fmt.Errorf("change supplier status: %w", err)
	}
	sup.Status = status
	if err := s.suppliers.Update(ctx, sup); err != nil {
		return nil, fmt.Errorf("change supplier status: %w", err)
	}
	return sup, nil
}

// Get 获取供应商详情。
func (s *SupplierService) Get(ctx context.Context, id uint) (*model.Supplier, error) {
	sup, err := s.suppliers.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeSupplierNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "供应商", id), err)
		}
		return nil, fmt.Errorf("get supplier: %w", err)
	}
	return sup, nil
}

// List 分页查询供应商。
func (s *SupplierService) List(ctx context.Context, page, pageSize int, status, category, name string) ([]model.Supplier, int64, error) {
	list, total, err := s.suppliers.List(ctx, page, pageSize, status, category, name)
	if err != nil {
		return nil, 0, fmt.Errorf("list suppliers: %w", err)
	}
	return list, total, nil
}

// ListInventory 查询某供应商的库存列表。
func (s *SupplierService) ListInventory(ctx context.Context, supplierID uint, page, pageSize int) ([]model.InventoryItem, int64, error) {
	list, total, err := s.inventory.ListBySupplier(ctx, supplierID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list supplier inventory: %w", err)
	}
	return list, total, nil
}

// ListOrders 查询某供应商的采购单列表。
func (s *SupplierService) ListOrders(ctx context.Context, supplierID uint, page, pageSize int) ([]model.PurchaseOrder, int64, error) {
	list, total, err := s.orders.List(ctx, page, pageSize, "", supplierID, "", "")
	if err != nil {
		return nil, 0, fmt.Errorf("list supplier orders: %w", err)
	}
	return list, total, nil
}

// CreateSupplierParams 新增供应商参数。
type CreateSupplierParams struct {
	Name          string
	ContactPerson string
	Phone         string
	Email         string
	Address       string
	Categories    []string
	Rating        float64
	Status        string
}

// UpdateSupplierParams 更新供应商参数。
type UpdateSupplierParams struct {
	Name          string
	ContactPerson string
	Phone         string
	Email         string
	Address       string
	Categories    []string
	Rating        float64
}
