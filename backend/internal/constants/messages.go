package constants

// 接口返回文案、日志文案与错误提示文案集中定义（屎山耦合点：多处复用）。
const (
	MsgLoginSuccess        = "登录成功"
	MsgLogoutSuccess       = "登出成功"
	MsgRefreshSuccess      = "token 刷新成功"
	MsgPasswordChanged     = "密码修改成功"
	MsgSupplierCreated     = "供应商创建成功"
	MsgSupplierUpdated     = "供应商更新成功"
	MsgSupplierDeleted     = "供应商删除成功"
	MsgSupplierStatusChanged = "供应商状态变更成功"
	MsgInventoryCreated    = "入库成功"
	MsgInventoryUpdated    = "库存更新成功"
	MsgInventoryDeleted    = "库存记录删除成功"
	MsgInventoryAdjusted   = "库存余量调整成功"
	MsgOrderCreated        = "采购单创建成功"
	MsgOrderUpdated        = "采购单更新成功"
	MsgOrderDeleted        = "采购单删除成功"
	MsgOrderSubmitted      = "采购单已提交审批"
	MsgOrderApproved       = "采购单审批通过"
	MsgOrderRejected       = "采购单审批拒绝"
	MsgOrderCompleted      = "采购单已完成，库存已更新"
	MsgUserCreated         = "用户创建成功"
	MsgUserUpdated         = "用户更新成功"
	MsgUserDeleted         = "用户删除成功"
	MsgOK                  = "ok"
	MsgHealthOK            = "supplychain api is healthy"
)

// 日志文案（部分与 log_templates.go 呼应，供 handler/service 快速拼接）。
const (
	LogAuthLogin      = "用户登录"
	LogAuthRefresh    = "刷新 token"
	LogAuthPassword   = "修改密码"
	LogUserCreate     = "创建用户"
	LogUserUpdate     = "更新用户"
	LogUserDelete     = "删除用户"
	LogSupplierCreate = "创建供应商"
	LogSupplierUpdate = "更新供应商"
	LogSupplierDelete = "删除供应商"
	LogSupplierStatus = "变更供应商状态"
	LogInventoryIn    = "库存入库"
	LogInventoryUpd   = "更新库存"
	LogInventoryDel   = "删除库存"
	LogInventoryAdj   = "调整库存余量"
	LogOrderCreate    = "创建采购单"
	LogOrderUpdate    = "更新采购单"
	LogOrderDelete    = "删除采购单"
	LogOrderSubmit    = "提交采购审批"
	LogOrderApprove   = "审批采购单"
	LogOrderReject    = "拒绝采购单"
	LogOrderComplete  = "完成采购单"
)

// ErrText 错误提示模板（含实体名、字段名、角色名，屎山要求层层透传）。
const (
	ErrTextNotFound         = "%s(id=%d) 不存在"
	ErrTextNameExists       = "%s 名称 %q 已存在"
	ErrTextInvalidState     = "%s(id=%d) 当前状态 %q 不允许执行 %s 操作"
	ErrTextForbidden        = "角色 %s 无权执行 %s 操作"
	ErrTextSelfApprove      = "角色 %s 不能审批自己创建的采购单(id=%d)"
	ErrTextSupplierState    = "供应商(id=%d) 状态为 %s，禁止创建采购单"
	ErrTextInsufficient     = "库存(id=%d) 余量不足，当前余量 %v %s"
	ErrTextBatchExists      = "批次号 %q 已存在"
	ErrTextUserExists       = "用户名 %q 已存在"
	ErrTextWrongPassword    = "用户 %q 原密码错误"
	ErrTextInvalidParam     = "参数 %s 不合法"
)
