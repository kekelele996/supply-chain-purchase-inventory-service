package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// SupplierHandler 供应商处理器。
type SupplierHandler struct {
	suppliers *service.SupplierService
	logger    *slog.Logger
}

// NewSupplierHandler 构造供应商处理器。
func NewSupplierHandler(suppliers *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{suppliers: suppliers, logger: util.Logger}
}

// List 供应商列表（分页 + 筛选）。
func (h *SupplierHandler) List(c *gin.Context) {
	page := util.ParsePageQuery(c)
	list, total, err := h.suppliers.List(c.Request.Context(), page.Page, page.PageSize,
		c.Query("status"), c.Query("category"), c.Query("name"))
	if err != nil {
		util.Error(c, err)
		return
	}
	items := make([]dto.SupplierResponse, 0, len(list))
	for _, s := range list {
		items = append(items, toSupplierResponse(s))
	}
	util.OK(c, dto.PageResult{List: items, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Create 新增供应商。
func (h *SupplierHandler) Create(c *gin.Context) {
	var req dto.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	sup, err := h.suppliers.Create(c.Request.Context(), service.CreateSupplierParams{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Categories:    req.Categories,
		Rating:        req.Rating,
		Status:        req.Status,
	})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Created(c, constants.MsgSupplierCreated, toSupplierResponse(*sup))
}

// Get 供应商详情。
func (h *SupplierHandler) Get(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	sup, err := h.suppliers.Get(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, toSupplierResponse(*sup))
}

// Update 更新供应商。
func (h *SupplierHandler) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	var req dto.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	sup, err := h.suppliers.Update(c.Request.Context(), idReq.ID, service.UpdateSupplierParams{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Categories:    req.Categories,
		Rating:        req.Rating,
	})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgSupplierUpdated, toSupplierResponse(*sup))
}

// Delete 删除供应商（软删除，admin）。
func (h *SupplierHandler) Delete(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	if err := h.suppliers.Delete(c.Request.Context(), idReq.ID); err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgSupplierDeleted, nil)
}

// ChangeStatus 变更供应商状态（manager+）。
func (h *SupplierHandler) ChangeStatus(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	var req dto.SupplierStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	sup, err := h.suppliers.ChangeStatus(c.Request.Context(), idReq.ID, constants.SupplierStatus(req.Status))
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgSupplierStatusChanged, toSupplierResponse(*sup))
}

// Inventory 查看某供应商的库存列表。
func (h *SupplierHandler) Inventory(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	page := util.ParsePageQuery(c)
	list, total, err := h.suppliers.ListInventory(c.Request.Context(), idReq.ID, page.Page, page.PageSize)
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

// Orders 查看某供应商的采购单列表。
func (h *SupplierHandler) Orders(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	page := util.ParsePageQuery(c)
	list, total, err := h.suppliers.ListOrders(c.Request.Context(), idReq.ID, page.Page, page.PageSize)
	if err != nil {
		util.Error(c, err)
		return
	}
	items := make([]dto.PurchaseOrderResponse, 0, len(list))
	for _, o := range list {
		items = append(items, toPurchaseOrderResponse(o))
	}
	util.OK(c, dto.PageResult{List: items, Total: total, Page: page.Page, PageSize: page.PageSize})
}
