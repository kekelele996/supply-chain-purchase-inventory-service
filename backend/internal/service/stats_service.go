package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

// StatsService 成本统计服务，直接聚合采购单与采购明细表。
type StatsService struct {
	db *gorm.DB
}

// NewStatsService 构造成本统计服务。
func NewStatsService(db *gorm.DB) *StatsService {
	return &StatsService{db: db}
}

// ByCategory 按食材分类统计采购成本（approved/completed 状态）。
func (s *StatsService) ByCategory(ctx context.Context, startDate, endDate string) ([]map[string]interface{}, error) {
	q := s.aggQuery(ctx, startDate, endDate).
		Joins("JOIN inventory_items inv ON inv.id = purchase_order_items.inventory_item_id").
		Select("inv.category AS category, ROUND(SUM(purchase_order_items.subtotal), 2) AS total_amount, COUNT(DISTINCT po.id) AS order_count").
		Group("inv.category").
		Order("total_amount DESC")
	return s.rows(q)
}

// BySupplier 按供应商统计采购总额。
func (s *StatsService) BySupplier(ctx context.Context, startDate, endDate string) ([]map[string]interface{}, error) {
	q := s.aggQuery(ctx, startDate, endDate).
		Select("po.supplier_id AS supplier_id, sup.name AS supplier_name, ROUND(SUM(po.total_amount), 2) AS total_amount, COUNT(DISTINCT po.id) AS order_count").
		Joins("JOIN suppliers sup ON sup.id = po.supplier_id").
		Group("po.supplier_id, sup.name").
		Order("total_amount DESC")
	return s.rows(q)
}

// Trend 按时间段统计成本趋势（period: day/week/month）。
func (s *StatsService) Trend(ctx context.Context, period, startDate, endDate string) ([]map[string]interface{}, error) {
	layout := "%Y-%m-%d"
	switch period {
	case "week":
		layout = "%Y-%u"
	case "month":
		layout = "%Y-%m"
	}
	q := s.aggQuery(ctx, startDate, endDate).
		Select("DATE_FORMAT(po.created_at, ?) AS period, ROUND(SUM(po.total_amount), 2) AS total_amount, COUNT(DISTINCT po.id) AS order_count", layout).
		Group("period").
		Order("period ASC")
	return s.rows(q)
}

// Summary 采购成本概览。
func (s *StatsService) Summary(ctx context.Context) (map[string]interface{}, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	base := s.db.WithContext(ctx).Model(&model.PurchaseOrder{}).
		Where("status IN ?", []string{string(constants.OrderApproved), string(constants.OrderCompleted)})

	var total float64
	var count int64
	if err := base.Select("COALESCE(SUM(total_amount),0)").Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("stats summary total: %w", err)
	}
	if err := base.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("stats summary count: %w", err)
	}
	var monthTotal float64
	var monthCount int64
	monthQ := base.Where("created_at >= ?", monthStart)
	if err := monthQ.Select("COALESCE(SUM(total_amount),0)").Scan(&monthTotal).Error; err != nil {
		return nil, fmt.Errorf("stats summary month total: %w", err)
	}
	if err := monthQ.Count(&monthCount).Error; err != nil {
		return nil, fmt.Errorf("stats summary month count: %w", err)
	}
	avg := 0.0
	if count > 0 {
		avg = round2(total / float64(count))
	}
	return map[string]interface{}{
		"total_purchase_amount": round2(total),
		"avg_order_amount":      avg,
		"month_purchase_amount": round2(monthTotal),
		"total_order_count":     count,
		"month_order_count":     monthCount,
	}, nil
}

func (s *StatsService) aggQuery(ctx context.Context, startDate, endDate string) *gorm.DB {
	q := s.db.WithContext(ctx).Model(&model.PurchaseOrderItem{}).
		Joins("JOIN purchase_orders po ON po.id = purchase_order_items.order_id").
		Where("po.status IN ?", []string{string(constants.OrderApproved), string(constants.OrderCompleted)})
	if startDate != "" {
		q = q.Where("po.created_at >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		q = q.Where("po.created_at <= ?", endDate+" 23:59:59")
	}
	return q
}

func (s *StatsService) rows(q *gorm.DB) ([]map[string]interface{}, error) {
	var out []map[string]interface{}
	if err := q.Scan(&out).Error; err != nil {
		return nil, fmt.Errorf("stats query: %w", err)
	}
	return out, nil
}
