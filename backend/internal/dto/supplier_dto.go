package dto

import "time"

// CreateSupplierRequest 新增供应商请求。
type CreateSupplierRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=100"`
	ContactPerson string   `json:"contact_person" binding:"required,min=1,max=50"`
	Phone         string   `json:"phone" binding:"required,min=5,max=20"`
	Email         string   `json:"email" binding:"omitempty,email,max=100"`
	Address       string   `json:"address" binding:"required,min=2,max=200"`
	Categories    []string `json:"categories" binding:"required,min=1,dive,min=1,max=30"`
	Rating        float64  `json:"rating" binding:"omitempty,gte=1.0,lte=5.0"`
	Status        string   `json:"status" binding:"omitempty,oneof=active suspended blacklisted"`
}

// UpdateSupplierRequest 更新供应商请求。
type UpdateSupplierRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=100"`
	ContactPerson string   `json:"contact_person" binding:"required,min=1,max=50"`
	Phone         string   `json:"phone" binding:"required,min=5,max=20"`
	Email         string   `json:"email" binding:"omitempty,email,max=100"`
	Address       string   `json:"address" binding:"required,min=2,max=200"`
	Categories    []string `json:"categories" binding:"required,min=1,dive,min=1,max=30"`
	Rating        float64  `json:"rating" binding:"omitempty,gte=1.0,lte=5.0"`
}

// SupplierStatusRequest 变更供应商状态请求。
type SupplierStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active suspended blacklisted"`
}

// SupplierResponse 供应商响应。
type SupplierResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	ContactPerson string    `json:"contact_person"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Address       string    `json:"address"`
	Categories    []string  `json:"categories"`
	Rating        float64   `json:"rating"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
