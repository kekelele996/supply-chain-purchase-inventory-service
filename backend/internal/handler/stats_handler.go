package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// StatsHandler 成本统计处理器。
type StatsHandler struct {
	stats  *service.StatsService
	logger *slog.Logger
}

// NewStatsHandler 构造成本统计处理器。
func NewStatsHandler(stats *service.StatsService) *StatsHandler {
	return &StatsHandler{stats: stats, logger: util.Logger}
}

// ByCategory 按食材分类统计采购成本。
func (h *StatsHandler) ByCategory(c *gin.Context) {
	data, err := h.stats.ByCategory(c.Request.Context(), c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, data)
}

// BySupplier 按供应商统计采购总额。
func (h *StatsHandler) BySupplier(c *gin.Context) {
	data, err := h.stats.BySupplier(c.Request.Context(), c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, data)
}

// Trend 成本趋势。
func (h *StatsHandler) Trend(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	data, err := h.stats.Trend(c.Request.Context(), period, c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, data)
}

// Summary 成本概览。
func (h *StatsHandler) Summary(c *gin.Context) {
	data, err := h.stats.Summary(c.Request.Context())
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, data)
}
