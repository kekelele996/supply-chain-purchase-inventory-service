package dto

import "time"

// PurchaseItemRequest 采购明细请求。
type PurchaseItemRequest struct {
	InventoryItemID uint    `json:"inventory_item_id" binding:"required,gt=0"`
	Quantity        float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice       float64 `json:"unit_price" binding:"required,gte=0"`
}

// CreatePurchaseRequest 创建采购单请求。
type CreatePurchaseRequest struct {
	SupplierID uint                   `json:"supplier_id" binding:"required,gt=0"`
	Notes      string                 `json:"notes" binding:"omitempty,max=500"`
	Items      []PurchaseItemRequest  `json:"items" binding:"required,min=1,dive"`
}

// UpdatePurchaseRequest 更新采购单请求（仅草稿/已拒绝状态）。
type UpdatePurchaseRequest struct {
	SupplierID uint                  `json:"supplier_id" binding:"required,gt=0"`
	Notes      string                `json:"notes" binding:"omitempty,max=500"`
	Items      []PurchaseItemRequest `json:"items" binding:"required,min=1,dive"`
}

// PurchaseOrderItemResponse 采购明细响应。
type PurchaseOrderItemResponse struct {
	ID              uint    `json:"id"`
	InventoryItemID uint    `json:"inventory_item_id"`
	ItemName        string  `json:"item_name"`
	Quantity        float64 `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	Subtotal        float64 `json:"subtotal"`
}

// PurchaseOrderResponse 采购单响应。
type PurchaseOrderResponse struct {
	ID          uint                       `json:"id"`
	OrderNo     string                     `json:"order_no"`
	SupplierID  uint                       `json:"supplier_id"`
	SupplierName string                    `json:"supplier_name"`
	Status      string                     `json:"status"`
	TotalAmount float64                    `json:"total_amount"`
	CreatorID   uint                       `json:"creator_id"`
	CreatorName string                     `json:"creator_name"`
	ApproverID  *uint                      `json:"approver_id"`
	ApproverName string                    `json:"approver_name"`
	ApprovedAt  *time.Time                 `json:"approved_at"`
	CompletedAt *time.Time                 `json:"completed_at"`
	Notes       string                     `json:"notes"`
	Items       []PurchaseOrderItemResponse `json:"items"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}
