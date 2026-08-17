package dto

import "time"

// OperationLogResponse 操作日志响应。
type OperationLogResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Action     string    `json:"action"`
	TargetType string    `json:"target_type"`
	TargetID   uint      `json:"target_id"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}
