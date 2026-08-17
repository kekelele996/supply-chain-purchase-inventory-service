package model

import (
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// PurchaseOrder 采购单。
type PurchaseOrder struct {
	ID          uint                          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo     string                        `gorm:"size:30;uniqueIndex;not null" json:"order_no"`
	SupplierID  uint                          `gorm:"not null;index" json:"supplier_id"`
	Supplier    *Supplier                     `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Status      constants.PurchaseOrderStatus `gorm:"type:varchar(20);not null;default:draft" json:"status"`
	TotalAmount float64                       `gorm:"type:decimal(12,2);not null;default:0" json:"total_amount"`
	CreatorID   uint                          `gorm:"not null;index" json:"creator_id"`
	Creator     *User                         `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	ApproverID  *uint                         `gorm:"index" json:"approver_id,omitempty"`
	Approver    *User                         `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
	ApprovedAt  *time.Time                    `json:"approved_at,omitempty"`
	CompletedAt *time.Time                    `json:"completed_at,omitempty"`
	Notes       string                        `gorm:"type:text" json:"notes"`
	Items       []PurchaseOrderItem           `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt   time.Time                     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time                     `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名。
func (PurchaseOrder) TableName() string { return "purchase_orders" }
