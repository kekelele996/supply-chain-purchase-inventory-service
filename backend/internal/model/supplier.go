package model

import (
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// Supplier 供应商。
type Supplier struct {
	ID            uint                       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string                     `gorm:"size:100;uniqueIndex;not null" json:"name"`
	ContactPerson string                     `gorm:"size:50;not null" json:"contact_person"`
	Phone         string                     `gorm:"size:20;not null" json:"phone"`
	Email         string                     `gorm:"size:100" json:"email"`
	Address       string                     `gorm:"size:200;not null" json:"address"`
	Categories    StringList                  `gorm:"type:json;not null" json:"categories"`
	Rating        float64                    `gorm:"not null;default:3.0" json:"rating"`
	Status        constants.SupplierStatus   `gorm:"type:varchar(20);not null;default:active" json:"status"`
	DeletedAt     *time.Time                 `gorm:"index" json:"-"`
	CreatedAt     time.Time                  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time                  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名。
func (Supplier) TableName() string { return "suppliers" }
