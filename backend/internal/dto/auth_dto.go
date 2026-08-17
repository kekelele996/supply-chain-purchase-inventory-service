package dto

import "time"

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// RefreshRequest 刷新 token 请求。
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest 修改密码请求。
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=72"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=72"`
}

// TokenResponse token 响应。
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserProfileResponse 当前用户信息响应。
type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
