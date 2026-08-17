package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// FormatDate 将时间格式化为 YYYY-MM-DD。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 将时间格式化为 YYYY-MM-DD HH:MM:SS。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseDate 解析 YYYY-MM-DD 日期，非法时返回零值与错误。
func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}

// SupplierStatusText 返回供应商状态的中文文案。
func SupplierStatusText(s constants.SupplierStatus) string {
	switch s {
	case constants.SupplierActive:
		return "活跃"
	case constants.SupplierSuspended:
		return "暂停"
	case constants.SupplierBlacklisted:
		return "黑名单"
	}
	return string(s)
}

// InventoryStatusText 返回库存状态的中文文案。
func InventoryStatusText(s constants.InventoryStatus) string {
	switch s {
	case constants.InventoryNormal:
		return "正常"
	case constants.InventoryLow:
		return "低于预警阈值"
	case constants.InventoryExpired:
		return "已过期"
	}
	return string(s)
}

// InventoryCategoryText 返回食材分类的中文文案。
func InventoryCategoryText(c constants.InventoryCategory) string {
	switch c {
	case constants.CategoryVegetable:
		return "蔬菜"
	case constants.CategoryMeat:
		return "肉类"
	case constants.CategorySeafood:
		return "海鲜"
	case constants.CategorySeasoning:
		return "调料"
	case constants.CategoryStaple:
		return "主食"
	case constants.CategoryDryGoods:
		return "干货"
	case constants.CategoryBeverage:
		return "饮品"
	}
	return string(c)
}

// OrderStatusText 返回采购单状态的中文文案。
func OrderStatusText(s constants.PurchaseOrderStatus) string {
	switch s {
	case constants.OrderDraft:
		return "草稿"
	case constants.OrderPendingApproval:
		return "待审批"
	case constants.OrderApproved:
		return "已批准"
	case constants.OrderRejected:
		return "已拒绝"
	case constants.OrderCompleted:
		return "已完成"
	}
	return string(s)
}

// RoleText 返回用户角色的中文文案。
func RoleText(r constants.UserRole) string {
	switch r {
	case constants.RoleOperator:
		return "操作员"
	case constants.RoleManager:
		return "经理"
	case constants.RoleAdmin:
		return "管理员"
	}
	return string(r)
}

// ActionText 返回操作动作的中文文案。
func ActionText(a constants.OperationAction) string {
	switch a {
	case constants.ActionCreate:
		return "创建"
	case constants.ActionUpdate:
		return "更新"
	case constants.ActionDelete:
		return "删除"
	case constants.ActionApprove:
		return "审批通过"
	case constants.ActionReject:
		return "审批拒绝"
	case constants.ActionLogin:
		return "登录"
	case constants.ActionLogout:
		return "登出"
	case constants.ActionSubmit:
		return "提交审批"
	case constants.ActionComplete:
		return "标记完成"
	}
	return string(a)
}

// Quantify 将数量与单位组合为可读文本。
func Quantify(q float64, unit string) string {
	return fmt.Sprintf("%g %s", q, strings.TrimSpace(unit))
}
