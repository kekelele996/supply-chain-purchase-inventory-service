// Package model 定义全部数据库实体。
package model

import (
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// User 系统用户。
type User struct {
	ID           uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string             `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string             `gorm:"size:200;not null" json:"-"`
	Role         constants.UserRole `gorm:"type:varchar(20);not null;default:operator" json:"role"`
	CreatedAt    time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名。
func (User) TableName() string { return "users" }
