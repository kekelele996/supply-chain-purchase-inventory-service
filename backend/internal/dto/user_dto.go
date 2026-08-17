package dto

import "time"

// CreateUserRequest 创建用户请求。
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Role     string `json:"role" binding:"required,oneof=operator manager admin"`
}

// UpdateUserRequest 更新用户请求。
type UpdateUserRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=72"`
	Role     string `json:"role" binding:"required,oneof=operator manager admin"`
}

// UserResponse 用户响应。
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
