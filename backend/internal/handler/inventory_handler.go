package handler

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// InventoryHandler 库存处理器。
type InventoryHandler struct {
	inventory *service.InventoryService
	logger    *slog.Logger
}

// NewInventoryHandler 构造库存处理器。
func NewInventoryHandler(inventory *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventory: inventory, logger: util.Logger}
}

// List 库存列表（分页 + 筛选）。
func (h *InventoryHandler) List(c *gin.Context) {
	page := util.ParsePageQuery(c)
	supplierID := parseUintQuery(c.Query("supplier_id"))
	list, total, err := h.inventory.List(c.Request.Context(), page.Page, page.PageSize,
		c.Query("status"), c.Query("category"), supplierID, c.Query("name"))
	if err != nil {
		util.Error(c, err)
		return
	}
	items := make([]dto.InventoryResponse, 0, len(list))
	for _, item := range list {
		items = append(items, toInventoryResponse(item))
	}
	util.OK(c, dto.PageResult{List: items, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Create 新增库存（入库）。
func (h *InventoryHandler) Create(c *gin.Context) {
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	item, err := h.inventory.Create(c.Request.Context(), service.CreateInventoryParams{
		Name:            req.Name,
		Category:        req.Category,
		SupplierID:      req.SupplierID,
		BatchNo:         req.BatchNo,
		Quantity:        req.Quantity,
		Unit:            req.Unit,
		MinThreshold:    req.MinThreshold,
		ExpiryDate:      req.ExpiryDate,
		StorageLocation: req.StorageLocation,
	})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Created(c, constants.MsgInventoryCreated, toInventoryResponse(*item))
}

// Get 库存详情。
func (h *InventoryHandler) Get(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	item, err := h.inventory.Get(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, toInventoryResponse(*item))
}

// Update 更新库存信息。
func (h *InventoryHandler) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	item, err := h.inventory.Update(c.Request.Context(), idReq.ID, service.UpdateInventoryParams{
		Name:            req.Name,
		Category:        req.Category,
		SupplierID:      req.SupplierID,
		Unit:            req.Unit,
		MinThreshold:    req.MinThreshold,
		ExpiryDate:      req.ExpiryDate,
		StorageLocation: req.StorageLocation,
	})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgInventoryUpdated, toInventoryResponse(*item))
}

// Delete 删除库存记录（admin）。
func (h *InventoryHandler) Delete(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	if err := h.inventory.Delete(c.Request.Context(), idReq.ID); err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgInventoryDeleted, nil)
}

// Alerts 库存预警列表。
func (h *InventoryHandler) Alerts(c *gin.Context) {
	list, err := h.inventory.Alerts(c.Request.Context())
	if err != nil {
		util.Error(c, err)
		return
	}
	items := make([]dto.InventoryResponse, 0, len(list))
	for _, item := range list {
		items = append(items, toInventoryResponse(item))
	}
	util.OK(c, items)
}

// AdjustQuantity 调整库存余量（出库/入库）。
func (h *InventoryHandler) AdjustQuantity(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	var req dto.QuantityAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	item, err := h.inventory.AdjustQuantity(c.Request.Context(), idReq.ID, req.Delta)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgInventoryAdjusted, toInventoryResponse(*item))
}

func parseUintQuery(s string) uint {
	if s == "" {
		return 0
	}
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}
