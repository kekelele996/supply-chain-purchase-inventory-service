// Package util 提供日志、错误、JWT、密码、格式化等通用工具。
package util

import (
	"context"
	"log/slog"
	"os"
)

// Logger 全局结构化日志实例（log/slog），供所有 handler/service/middleware 引用。
var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

// SetLogger 替换全局日志实例（测试或自定义格式时使用）。
func SetLogger(l *slog.Logger) {
	Logger = l
}

// L 返回带请求上下文字段的日志器。
func L(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Logger
	}
	reqID := RequestIDFromContext(ctx)
	if reqID == "" {
		return Logger
	}
	return Logger.With("request_id", reqID)
}
