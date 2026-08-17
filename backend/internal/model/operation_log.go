package model

import (
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// OperationLog 操作日志。
type OperationLog struct {
	ID         uint                         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint                         `gorm:"not null;index" json:"user_id"`
	User       *User                        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     constants.OperationAction    `gorm:"type:varchar(20);not null" json:"action"`
	TargetType string                       `gorm:"size:50;not null" json:"target_type"`
	TargetID   uint                         `gorm:"not null" json:"target_id"`
	Detail     string                       `gorm:"type:json" json:"detail"`
	CreatedAt  time.Time                    `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名。
func (OperationLog) TableName() string { return "operation_logs" }
