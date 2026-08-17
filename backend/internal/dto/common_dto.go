// Package dto 定义请求与响应数据结构。
package dto

// PageResult 统一分页返回结构。
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// IDRequest 通用 ID 路径参数。
type IDRequest struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}
