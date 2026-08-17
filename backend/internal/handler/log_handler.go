package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// LogHandler 操作日志处理器。
type LogHandler struct {
	logs   *service.OperationLogService
	logger *slog.Logger
}

// NewLogHandler 构造操作日志处理器。
func NewLogHandler(logs *service.OperationLogService) *LogHandler {
	return &LogHandler{logs: logs, logger: util.Logger}
}

// List 操作日志列表（分页 + 筛选）。
func (h *LogHandler) List(c *gin.Context) {
	page := util.ParsePageQuery(c)
	userID := parseUintQuery(c.Query("user_id"))
	list, total, err := h.logs.List(c.Request.Context(), page.Page, page.PageSize,
		userID, c.Query("action"), c.Query("target_type"), c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		util.Error(c, err)
		return
	}
	items := make([]dto.OperationLogResponse, 0, len(list))
	for _, l := range list {
		items = append(items, toOperationLogResponse(l))
	}
	util.OK(c, dto.PageResult{List: items, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Get 日志详情。
func (h *LogHandler) Get(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	entry, err := h.logs.Get(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, toOperationLogResponse(*entry))
}
