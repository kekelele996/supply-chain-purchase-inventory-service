package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageQuery 统一分页参数。
type PageQuery struct {
	Page     int
	PageSize int
}

// ParsePageQuery 从查询参数解析 page/page_size，带默认值与上下限约束。
func ParsePageQuery(c *gin.Context) PageQuery {
	page := parsePositiveInt(c.Query("page"), 1)
	pageSize := parsePositiveInt(c.Query("page_size"), 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PageQuery{Page: page, PageSize: pageSize}
}

// Offset 返回 SQL offset。
func (p PageQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 返回 SQL limit。
func (p PageQuery) Limit() int {
	return p.PageSize
}

func parsePositiveInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
