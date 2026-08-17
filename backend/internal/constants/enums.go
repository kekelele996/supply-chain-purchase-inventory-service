// Package constants 集中维护业务枚举、错误码、文案与日志模板。
package constants

import "strings"

// UserRole 用户角色枚举。
type UserRole string

const (
	RoleOperator UserRole = "operator"
	RoleManager  UserRole = "manager"
	RoleAdmin    UserRole = "admin"
)

// ValidRoles 返回全部合法角色。
func ValidRoles() []UserRole { return []UserRole{RoleOperator, RoleManager, RoleAdmin} }

// IsValidUserRole 校验角色是否合法。
func IsValidUserRole(r UserRole) bool {
	switch r {
	case RoleOperator, RoleManager, RoleAdmin:
		return true
	}
	return false
}

// SupplierStatus 供应商状态枚举。
type SupplierStatus string

const (
	SupplierActive      SupplierStatus = "active"
	SupplierSuspended   SupplierStatus = "suspended"
	SupplierBlacklisted SupplierStatus = "blacklisted"
)

// ValidSupplierStatuses 返回全部合法供应商状态。
func ValidSupplierStatuses() []SupplierStatus {
	return []SupplierStatus{SupplierActive, SupplierSuspended, SupplierBlacklisted}
}

// IsValidSupplierStatus 校验供应商状态。
func IsValidSupplierStatus(s SupplierStatus) bool {
	switch s {
	case SupplierActive, SupplierSuspended, SupplierBlacklisted:
		return true
	}
	return false
}

// InventoryCategory 食材分类枚举。
type InventoryCategory string

const (
	CategoryVegetable InventoryCategory = "vegetable"
	CategoryMeat      InventoryCategory = "meat"
	CategorySeafood   InventoryCategory = "seafood"
	CategorySeasoning InventoryCategory = "seasoning"
	CategoryStaple    InventoryCategory = "staple"
	CategoryDryGoods  InventoryCategory = "dry_goods"
	CategoryBeverage  InventoryCategory = "beverage"
)

// ValidInventoryCategories 返回全部合法食材分类。
func ValidInventoryCategories() []InventoryCategory {
	return []InventoryCategory{CategoryVegetable, CategoryMeat, CategorySeafood, CategorySeasoning, CategoryStaple, CategoryDryGoods, CategoryBeverage}
}

// IsValidInventoryCategory 校验食材分类。
func IsValidInventoryCategory(c InventoryCategory) bool {
	switch c {
	case CategoryVegetable, CategoryMeat, CategorySeafood, CategorySeasoning, CategoryStaple, CategoryDryGoods, CategoryBeverage:
		return true
	}
	return false
}

// InventoryStatus 库存状态枚举。
type InventoryStatus string

const (
	InventoryNormal  InventoryStatus = "normal"
	InventoryLow     InventoryStatus = "low"
	InventoryExpired InventoryStatus = "expired"
)

// ValidInventoryStatuses 返回全部合法库存状态。
func ValidInventoryStatuses() []InventoryStatus {
	return []InventoryStatus{InventoryNormal, InventoryLow, InventoryExpired}
}

// IsValidInventoryStatus 校验库存状态。
func IsValidInventoryStatus(s InventoryStatus) bool {
	switch s {
	case InventoryNormal, InventoryLow, InventoryExpired:
		return true
	}
	return false
}

// PurchaseOrderStatus 采购单状态枚举。
type PurchaseOrderStatus string

const (
	OrderDraft          PurchaseOrderStatus = "draft"
	OrderPendingApproval PurchaseOrderStatus = "pending_approval"
	OrderApproved       PurchaseOrderStatus = "approved"
	OrderRejected       PurchaseOrderStatus = "rejected"
	OrderCompleted      PurchaseOrderStatus = "completed"
)

// ValidPurchaseOrderStatuses 返回全部合法采购单状态。
func ValidPurchaseOrderStatuses() []PurchaseOrderStatus {
	return []PurchaseOrderStatus{OrderDraft, OrderPendingApproval, OrderApproved, OrderRejected, OrderCompleted}
}

// IsValidPurchaseOrderStatus 校验采购单状态。
func IsValidPurchaseOrderStatus(s PurchaseOrderStatus) bool {
	switch s {
	case OrderDraft, OrderPendingApproval, OrderApproved, OrderRejected:
		return true
	}
	return false
}

// OperationAction 操作日志动作枚举。
type OperationAction string

const (
	ActionCreate OperationAction = "create"
	ActionUpdate OperationAction = "update"
	ActionDelete OperationAction = "delete"
	ActionApprove OperationAction = "approve"
	ActionReject  OperationAction = "reject"
	ActionLogin   OperationAction = "login"
	ActionLogout  OperationAction = "logout"
	ActionSubmit  OperationAction = "submit"
	ActionComplete OperationAction = "complete"
)

// ValidOperationActions 返回全部合法操作动作。
func ValidOperationActions() []OperationAction {
	return []OperationAction{ActionCreate, ActionUpdate, ActionDelete, ActionApprove, ActionReject, ActionLogin, ActionLogout, ActionSubmit, ActionComplete}
}

// IsValidOperationAction 校验操作动作。
func IsValidOperationAction(a OperationAction) bool {
	switch a {
	case ActionCreate, ActionUpdate, ActionDelete, ActionApprove, ActionReject, ActionLogin, ActionLogout, ActionSubmit, ActionComplete:
		return true
	}
	return false
}

// NormalizeRole 将字符串转换为 UserRole，非法值返回空。
func NormalizeRole(s string) UserRole {
	r := UserRole(strings.ToLower(strings.TrimSpace(s)))
	if IsValidUserRole(r) {
		return r
	}
	return ""
}
