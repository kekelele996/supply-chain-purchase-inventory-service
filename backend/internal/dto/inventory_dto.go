package dto

import "time"

// CreateInventoryRequest 新增库存（入库）请求。
type CreateInventoryRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=100"`
	Category       string  `json:"category" binding:"required,oneof=vegetable meat seafood seasoning staple dry_goods beverage"`
	SupplierID     uint    `json:"supplier_id" binding:"required,gt=0"`
	BatchNo        string  `json:"batch_no" binding:"required,min=1,max=50"`
	Quantity       float64 `json:"quantity" binding:"required,gte=0"`
	Unit           string  `json:"unit" binding:"required,min=1,max=20"`
	MinThreshold   float64 `json:"min_threshold" binding:"omitempty,gte=0"`
	ExpiryDate     string  `json:"expiry_date" binding:"required,datetime=2006-01-02"`
	StorageLocation string `json:"storage_location" binding:"required,min=1,max=100"`
}

// UpdateInventoryRequest 更新库存信息请求。
type UpdateInventoryRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=100"`
	Category       string  `json:"category" binding:"required,oneof=vegetable meat seafood seasoning staple dry_goods beverage"`
	SupplierID     uint    `json:"supplier_id" binding:"required,gt=0"`
	Unit           string  `json:"unit" binding:"required,min=1,max=20"`
	MinThreshold   float64 `json:"min_threshold" binding:"omitempty,gte=0"`
	ExpiryDate     string  `json:"expiry_date" binding:"required,datetime=2006-01-02"`
	StorageLocation string `json:"storage_location" binding:"required,min=1,max=100"`
}

// QuantityAdjustRequest 调整库存余量请求。
type QuantityAdjustRequest struct {
	Delta float64 `json:"delta" binding:"required"`
	Reason string `json:"reason" binding:"omitempty,max=200"`
}

// InventoryResponse 库存响应。
type InventoryResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	SupplierID      uint      `json:"supplier_id"`
	SupplierName    string    `json:"supplier_name"`
	BatchNo         string    `json:"batch_no"`
	Quantity        float64   `json:"quantity"`
	Unit            string    `json:"unit"`
	MinThreshold    float64   `json:"min_threshold"`
	ExpiryDate      string    `json:"expiry_date"`
	StorageLocation string    `json:"storage_location"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
