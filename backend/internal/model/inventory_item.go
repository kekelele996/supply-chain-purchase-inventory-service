package model

import (
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// InventoryItem 食材库存。
type InventoryItem struct {
	ID             uint                        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string                      `gorm:"size:100;not null" json:"name"`
	Category       constants.InventoryCategory `gorm:"type:varchar(20);not null" json:"category"`
	SupplierID     uint                        `gorm:"not null;index" json:"supplier_id"`
	Supplier       *Supplier                   `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	BatchNo        string                      `gorm:"size:50;uniqueIndex;not null" json:"batch_no"`
	Quantity       float64                     `gorm:"not null;default:0" json:"quantity"`
	Unit           string                      `gorm:"size:20;not null" json:"unit"`
	MinThreshold   float64                     `gorm:"not null;default:10.0" json:"min_threshold"`
	ExpiryDate     time.Time                   `gorm:"type:date;not null" json:"expiry_date"`
	StorageLocation string                     `gorm:"size:100;not null" json:"storage_location"`
	Status         constants.InventoryStatus   `gorm:"type:varchar(20);not null;default:normal" json:"status"`
	CreatedAt      time.Time                   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time                   `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名。
func (InventoryItem) TableName() string { return "inventory_items" }
