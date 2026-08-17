package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// PurchaseHandler 采购单处理器。
type PurchaseHandler struct {
	purchases *service.PurchaseService
	logger    *slog.Logger
}

// NewPurchaseHandler 构造采购单处理器。
func NewPurchaseHandler(purchases *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchases: purchases, logger: util.Logger}
}

// List 采购单列表（分页 + 筛选）。
func (h *PurchaseHandler) List(c *gin.Context) {
	page := util.ParsePageQuery(c)
	supplierID := parseUintQuery(c.Query("supplier_id"))
	list, total, err := h.purchases.List(c.Request.Context(), page.Page, page.PageSize,
		c.Query("status"), supplierID, c.Query("start_date"), c.Query("end_date"))
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

// Create 创建采购单（含明细）。
func (h *PurchaseHandler) Create(c *gin.Context) {
	var req dto.CreatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user := util.CurrentUser(c)
	order, err := h.purchases.Create(c.Request.Context(), service.CreatePurchaseParams{
		SupplierID: req.SupplierID,
		Notes:      req.Notes,
		Items:      toItemParams(req.Items),
		CreatorID:  user.ID,
	})
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Created(c, constants.MsgOrderCreated, toPurchaseOrderResponse(*order))
}

// Get 采购单详情（含明细）。
func (h *PurchaseHandler) Get(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	order, err := h.purchases.Get(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, toPurchaseOrderResponse(*order))
}

// Update 更新采购单（仅草稿/已拒绝）。
func (h *PurchaseHandler) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	var req dto.UpdatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user := util.CurrentUser(c)
	order, err := h.purchases.Update(c.Request.Context(), idReq.ID, service.UpdatePurchaseParams{
		SupplierID: req.SupplierID,
		Notes:      req.Notes,
		Items:      toItemParams(req.Items),
	}, user.ID, user.Role)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOrderUpdated, toPurchaseOrderResponse(*order))
}

// Delete 删除采购单（仅草稿）。
func (h *PurchaseHandler) Delete(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user := util.CurrentUser(c)
	if err := h.purchases.Delete(c.Request.Context(), idReq.ID, user.ID, user.Role); err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOrderDeleted, nil)
}

// Submit 提交审批。
func (h *PurchaseHandler) Submit(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	order, err := h.purchases.Submit(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOrderSubmitted, toPurchaseOrderResponse(*order))
}

// Approve 审批通过。
func (h *PurchaseHandler) Approve(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user := util.CurrentUser(c)
	order, err := h.purchases.Approve(c.Request.Context(), idReq.ID, user.ID, user.Role)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOrderApproved, toPurchaseOrderResponse(*order))
}

// Reject 审批拒绝。
func (h *PurchaseHandler) Reject(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user := util.CurrentUser(c)
	order, err := h.purchases.Reject(c.Request.Context(), idReq.ID, user.ID, user.Role)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOrderRejected, toPurchaseOrderResponse(*order))
}

// Complete 标记完成（自动更新库存）。
func (h *PurchaseHandler) Complete(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	order, items, err := h.purchases.Complete(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOrderCompleted, gin.H{
		"order":          toPurchaseOrderResponse(*order),
		"items_updated":  items,
	})
}

func toItemParams(items []dto.PurchaseItemRequest) []service.PurchaseItemParams {
	out := make([]service.PurchaseItemParams, 0, len(items))
	for _, it := range items {
		out = append(out, service.PurchaseItemParams{
			InventoryItemID: it.InventoryItemID,
			Quantity:        it.Quantity,
			UnitPrice:       it.UnitPrice,
		})
	}
	return out
}
