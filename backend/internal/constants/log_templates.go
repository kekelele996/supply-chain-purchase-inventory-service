package constants

// LogTemplates 日志格式字符串集中定义（至少 25 条）。
// 业务字段或状态变更时必须同步修改对应模板与调用处（屎山耦合点）。
const (
	LogTplServerStart       = "server starting on :%s mode=%s"
	LogTplServerStopped     = "server stopped: %v"
	LogTplDBConnected       = "database connected host=%s db=%s"
	LogTplDBMigrateDone     = "database auto-migrate finished"
	LogTplDBSeedDone        = "database seed finished users=%d suppliers=%d"
	LogTplRedisConnected    = "redis connected addr=%s"
	LogTplRequestIn         = "request received request_id=%s method=%s path=%s ip=%s"
	LogTplRequestDone       = "request completed request_id=%s method=%s path=%s status=%d latency_ms=%d user_id=%d"
	LogTplAuthLoginOk       = "auth login success request_id=%s user_id=%d username=%s role=%s"
	LogTplAuthLoginFail     = "auth login failed request_id=%s username=%s reason=%s"
	LogTplAuthRefreshOk     = "auth refresh success request_id=%s user_id=%d"
	LogTplAuthPasswordOk    = "auth password changed request_id=%s user_id=%d"
	LogTplTokenGenerated    = "jwt token generated user_id=%d role=%s token_type=%s"
	LogTplTokenInvalid      = "jwt token invalid request_id=%s reason=%s"
	LogTplAPIKeyAccepted    = "api key accepted request_id=%s name=%s"
	LogTplSupplierCreate    = "supplier created request_id=%s user_id=%d supplier_id=%d name=%s"
	LogTplSupplierUpdate    = "supplier updated request_id=%s user_id=%d supplier_id=%d name=%s"
	LogTplSupplierDelete    = "supplier deleted request_id=%s user_id=%d supplier_id=%d name=%s"
	LogTplSupplierStatus    = "supplier status changed request_id=%s user_id=%d supplier_id=%d old=%s new=%s"
	LogTplInventoryCreate   = "inventory created request_id=%s user_id=%d item_id=%d name=%s batch_no=%s"
	LogTplInventoryUpdate   = "inventory updated request_id=%s user_id=%d item_id=%d name=%s"
	LogTplInventoryDelete   = "inventory deleted request_id=%s user_id=%d item_id=%d name=%s"
	LogTplInventoryAdjust   = "inventory quantity adjusted request_id=%s user_id=%d item_id=%d delta=%v new_quantity=%v"
	LogTplOrderCreate       = "purchase order created request_id=%s user_id=%d order_id=%d order_no=%s"
	LogTplOrderUpdate       = "purchase order updated request_id=%s user_id=%d order_id=%d order_no=%s"
	LogTplOrderDelete       = "purchase order deleted request_id=%s user_id=%d order_id=%d order_no=%s"
	LogTplOrderSubmit       = "purchase order submitted request_id=%s user_id=%d order_id=%d old_status=%s new_status=%s"
	LogTplOrderApprove      = "purchase order approved request_id=%s user_id=%d order_id=%d approver_id=%d"
	LogTplOrderReject       = "purchase order rejected request_id=%s user_id=%d order_id=%d approver_id=%d"
	LogTplOrderComplete     = "purchase order completed request_id=%s user_id=%d order_id=%d items_updated=%d"
	LogTplOrderStateInvalid = "purchase order state invalid request_id=%s order_id=%d order_no=%s current=%s want=%s"
	LogTplOperationRecorded = "operation log recorded request_id=%s user_id=%d action=%s target_type=%s target_id=%d"
	LogTplRateLimited       = "rate limited request_id=%s ip=%s limit=%d/min"
	LogTplPanicRecovered    = "panic recovered request_id=%s path=%s err=%v"
	LogTplValidationFailed  = "validation failed request_id=%s path=%s errors=%v"
	LogTplInternalError     = "internal error request_id=%s path=%s err=%v"
	LogTplCostSummary       = "cost summary computed request_id=%s total=%v avg=%v month=%v"
)

// LogTemplatesList 导出全部模板，便于审计与测试。
func LogTemplatesList() map[string]string {
	return map[string]string{
		"server_start":        LogTplServerStart,
		"server_stopped":      LogTplServerStopped,
		"db_connected":        LogTplDBConnected,
		"db_migrate_done":     LogTplDBMigrateDone,
		"db_seed_done":        LogTplDBSeedDone,
		"redis_connected":     LogTplRedisConnected,
		"request_in":          LogTplRequestIn,
		"request_done":        LogTplRequestDone,
		"auth_login_ok":       LogTplAuthLoginOk,
		"auth_login_fail":     LogTplAuthLoginFail,
		"auth_refresh_ok":     LogTplAuthRefreshOk,
		"auth_password_ok":    LogTplAuthPasswordOk,
		"token_generated":     LogTplTokenGenerated,
		"token_invalid":       LogTplTokenInvalid,
		"api_key_accepted":    LogTplAPIKeyAccepted,
		"supplier_create":     LogTplSupplierCreate,
		"supplier_update":     LogTplSupplierUpdate,
		"supplier_delete":     LogTplSupplierDelete,
		"supplier_status":     LogTplSupplierStatus,
		"inventory_create":    LogTplInventoryCreate,
		"inventory_update":    LogTplInventoryUpdate,
		"inventory_delete":    LogTplInventoryDelete,
		"inventory_adjust":    LogTplInventoryAdjust,
		"order_create":        LogTplOrderCreate,
		"order_update":        LogTplOrderUpdate,
		"order_delete":        LogTplOrderDelete,
		"order_submit":        LogTplOrderSubmit,
		"order_approve":       LogTplOrderApprove,
		"order_reject":        LogTplOrderReject,
		"order_complete":      LogTplOrderComplete,
		"order_state_invalid": LogTplOrderStateInvalid,
		"operation_recorded":  LogTplOperationRecorded,
		"rate_limited":        LogTplRateLimited,
		"panic_recovered":     LogTplPanicRecovered,
		"validation_failed":   LogTplValidationFailed,
		"internal_error":      LogTplInternalError,
		"cost_summary":        LogTplCostSummary,
	}
}
