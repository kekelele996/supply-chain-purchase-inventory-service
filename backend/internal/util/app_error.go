package util

import (
	"errors"
	"fmt"
)

// AppError 业务异常，携带统一错误码、HTTP 状态码与用户可读信息。
type AppError struct {
	Code    int
	HTTP    int
	Message string
	Err     error
}

// Error 实现 error 接口，保留错误链。
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d http=%d msg=%s cause=%v", e.Code, e.HTTP, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d http=%d msg=%s", e.Code, e.HTTP, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 透传底层错误。
func (e *AppError) Unwrap() error { return e.Err }

// NewAppError 构造业务异常。
func NewAppError(code, http int, message string, cause error) *AppError {
	return &AppError{Code: code, HTTP: http, Message: message, Err: cause}
}

// WrapAppError 将普通错误包装为业务异常。
func WrapAppError(code, http int, message string, cause error) error {
	return NewAppError(code, http, message, cause)
}

// AsAppError 将任意 error 转换为 *AppError；无法转换时返回 nil。
func AsAppError(err error) *AppError {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return nil
}
