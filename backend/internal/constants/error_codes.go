package constants

// 业务错误码集中定义。规则：0 表示成功，1xxx 通用，2xxx 认证授权，
// 3xxx 供应商，4xxx 库存，5xxx 采购，6xxx 用户，7xxx 日志/统计。
const (
	CodeOK                  = 0
	CodeBadRequest          = 1001
	CodeValidationError     = 1002
	CodeNotFound            = 1003
	CodeConflict            = 1004
	CodeInternalError       = 1005
	CodeUnauthorized        = 2001
	CodeForbidden           = 2002
	CodeTokenExpired        = 2003
	CodeTokenInvalid        = 2004
	CodeLoginFailed         = 2005
	CodePasswordWrong       = 2006
	CodeRateLimited         = 2007
	CodeAPIKeyInvalid       = 2008
	CodeSupplierNotFound    = 3001
	CodeSupplierNameExists  = 3002
	CodeSupplierForbidden   = 3003
	CodeInventoryNotFound   = 4001
	CodeInventoryBatchExists = 4002
	CodeInventoryInsufficient = 4003
	CodeInventoryExpired    = 4004
	CodeOrderNotFound       = 5001
	CodeOrderStateInvalid   = 5002
	CodeOrderCannotSelfApprove = 5003
	CodeOrderSupplierSuspended = 5004
	CodeUserNotFound        = 6001
	CodeUserNameExists      = 6002
	CodeLogNotFound         = 7001
)

// ErrorMessages 错误码到默认文案的映射。
var ErrorMessages = map[int]string{
	CodeOK:                    "ok",
	CodeBadRequest:            "请求参数错误",
	CodeValidationError:       "参数校验失败",
	CodeNotFound:              "资源不存在",
	CodeConflict:              "资源冲突",
	CodeInternalError:         "系统内部错误",
	CodeUnauthorized:          "未认证或登录已过期",
	CodeForbidden:             "权限不足",
	CodeTokenExpired:          "token 已过期",
	CodeTokenInvalid:          "token 无效",
	CodeLoginFailed:           "用户名或密码错误",
	CodePasswordWrong:         "原密码错误",
	CodeRateLimited:           "请求过于频繁，请稍后再试",
	CodeAPIKeyInvalid:         "API Key 无效",
	CodeSupplierNotFound:      "供应商不存在",
	CodeSupplierNameExists:    "供应商名称已存在",
	CodeSupplierForbidden:     "供应商状态不允许该操作",
	CodeInventoryNotFound:     "库存记录不存在",
	CodeInventoryBatchExists:  "批次号已存在",
	CodeInventoryInsufficient: "库存余量不足",
	CodeInventoryExpired:      "库存已过期，禁止操作",
	CodeOrderNotFound:         "采购单不存在",
	CodeOrderStateInvalid:     "采购单当前状态不允许该操作",
	CodeOrderCannotSelfApprove: "不能审批自己创建的采购单",
	CodeOrderSupplierSuspended: "供应商处于暂停或黑名单状态，无法创建采购单",
	CodeUserNotFound:          "用户不存在",
	CodeUserNameExists:        "用户名已存在",
	CodeLogNotFound:           "日志记录不存在",
}

// MessageOf 返回错误码对应文案，未登记时返回通用内部错误文案。
func MessageOf(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return ErrorMessages[CodeInternalError]
}
