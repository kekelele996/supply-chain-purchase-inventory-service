package dto

// CostByCategoryItem 按分类成本统计项。
type CostByCategoryItem struct {
	Category    string  `json:"category"`
	CategoryText string `json:"category_text"`
	TotalAmount float64 `json:"total_amount"`
	OrderCount  int64   `json:"order_count"`
}

// CostBySupplierItem 按供应商成本统计项。
type CostBySupplierItem struct {
	SupplierID   uint    `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	TotalAmount  float64 `json:"total_amount"`
	OrderCount   int64   `json:"order_count"`
}

// CostTrendItem 成本趋势项。
type CostTrendItem struct {
	Period      string  `json:"period"`
	TotalAmount float64 `json:"total_amount"`
	OrderCount  int64   `json:"order_count"`
}

// CostSummaryResponse 成本概览响应。
type CostSummaryResponse struct {
	TotalPurchaseAmount float64 `json:"total_purchase_amount"`
	AvgOrderAmount      float64 `json:"avg_order_amount"`
	MonthPurchaseAmount float64 `json:"month_purchase_amount"`
	TotalOrderCount     int64   `json:"total_order_count"`
	MonthOrderCount     int64   `json:"month_order_count"`
}
