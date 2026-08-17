package model

// PurchaseOrderItem 采购单明细。
type PurchaseOrderItem struct {
	ID              uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID         uint     `gorm:"not null;index" json:"order_id"`
	InventoryItemID uint     `gorm:"not null;index" json:"inventory_item_id"`
	InventoryItem   *InventoryItem `gorm:"foreignKey:InventoryItemID" json:"inventory_item,omitempty"`
	Quantity        float64  `gorm:"not null" json:"quantity"`
	UnitPrice       float64  `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal        float64  `gorm:"type:decimal(12,2);not null;default:0" json:"subtotal"`
}

// TableName 指定表名。
func (PurchaseOrderItem) TableName() string { return "purchase_order_items" }
