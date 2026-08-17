// Package docs 提供 OpenAPI 2.0 规范与 Swagger UI 数据。
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "餐饮供应链 API",
    "description": "餐饮供应链后端管理系统：供应商、库存、采购、成本统计与操作日志。",
    "version": "1.0.0"
  },
  "host": "{{.Host}}",
  "basePath": "/api/v1",
  "schemes": ["http"],
  "securityDefinitions": {
    "BearerAuth": {"type": "apiKey", "name": "Authorization", "in": "header", "description": "格式：Bearer <JWT>"},
    "ApiKeyAuth": {"type": "apiKey", "name": "X-API-Key", "in": "header"}
  },
  "tags": [
    {"name": "auth", "description": "认证模块"},
    {"name": "users", "description": "用户管理（admin）"},
    {"name": "suppliers", "description": "供应商管理"},
    {"name": "inventory", "description": "库存管理"},
    {"name": "purchases", "description": "采购管理"},
    {"name": "stats", "description": "成本统计"},
    {"name": "logs", "description": "操作日志"}
  ],
  "paths": {
    "/auth/login": {"post": {"tags": ["auth"], "summary": "用户登录", "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/LoginRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/auth/refresh": {"post": {"tags": ["auth"], "summary": "刷新 token", "security": [{"BearerAuth": []}], "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/RefreshRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/auth/me": {"get": {"tags": ["auth"], "summary": "获取当前用户信息", "security": [{"BearerAuth": []}], "responses": {"200": {"description": "OK"}}}},
    "/auth/password": {"put": {"tags": ["auth"], "summary": "修改密码", "security": [{"BearerAuth": []}], "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/ChangePasswordRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/users": {"get": {"tags": ["users"], "summary": "用户列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}, {"in": "query", "name": "username", "type": "string"}, {"in": "query", "name": "role", "type": "string"}], "responses": {"200": {"description": "OK"}}}, "post": {"tags": ["users"], "summary": "创建用户", "security": [{"BearerAuth": []}], "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/CreateUserRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/users/{id}": {"get": {"tags": ["users"], "summary": "用户详情", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}, "put": {"tags": ["users"], "summary": "更新用户", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateUserRequest"}}], "responses": {"200": {"description": "OK"}}}, "delete": {"tags": ["users"], "summary": "删除用户", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/suppliers": {"get": {"tags": ["suppliers"], "summary": "供应商列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}, {"in": "query", "name": "status", "type": "string"}, {"in": "query", "name": "category", "type": "string"}, {"in": "query", "name": "name", "type": "string"}], "responses": {"200": {"description": "OK"}}}, "post": {"tags": ["suppliers"], "summary": "新增供应商", "security": [{"BearerAuth": []}], "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/CreateSupplierRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/suppliers/{id}": {"get": {"tags": ["suppliers"], "summary": "供应商详情", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}, "put": {"tags": ["suppliers"], "summary": "更新供应商", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateSupplierRequest"}}], "responses": {"200": {"description": "OK"}}}, "delete": {"tags": ["suppliers"], "summary": "删除供应商（软删除，admin）", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/suppliers/{id}/status": {"put": {"tags": ["suppliers"], "summary": "变更供应商状态", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/SupplierStatusRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/suppliers/{id}/inventory": {"get": {"tags": ["suppliers"], "summary": "供应商库存列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/suppliers/{id}/orders": {"get": {"tags": ["suppliers"], "summary": "供应商采购单列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/inventory": {"get": {"tags": ["inventory"], "summary": "库存列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}, {"in": "query", "name": "status", "type": "string"}, {"in": "query", "name": "category", "type": "string"}, {"in": "query", "name": "supplier_id", "type": "integer"}, {"in": "query", "name": "name", "type": "string"}], "responses": {"200": {"description": "OK"}}}, "post": {"tags": ["inventory"], "summary": "新增库存（入库）", "security": [{"BearerAuth": []}], "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/CreateInventoryRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/inventory/alerts": {"get": {"tags": ["inventory"], "summary": "库存预警列表", "security": [{"BearerAuth": []}], "responses": {"200": {"description": "OK"}}}},
    "/inventory/{id}": {"get": {"tags": ["inventory"], "summary": "库存详情", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}, "put": {"tags": ["inventory"], "summary": "更新库存信息", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateInventoryRequest"}}], "responses": {"200": {"description": "OK"}}}, "delete": {"tags": ["inventory"], "summary": "删除库存记录（admin）", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/inventory/{id}/quantity": {"put": {"tags": ["inventory"], "summary": "调整库存余量", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/QuantityAdjustRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/purchases": {"get": {"tags": ["purchases"], "summary": "采购单列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}, {"in": "query", "name": "status", "type": "string"}, {"in": "query", "name": "supplier_id", "type": "integer"}, {"in": "query", "name": "start_date", "type": "string"}, {"in": "query", "name": "end_date", "type": "string"}], "responses": {"200": {"description": "OK"}}}, "post": {"tags": ["purchases"], "summary": "创建采购单", "security": [{"BearerAuth": []}], "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/CreatePurchaseRequest"}}], "responses": {"200": {"description": "OK"}}}},
    "/purchases/{id}": {"get": {"tags": ["purchases"], "summary": "采购单详情", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}, "put": {"tags": ["purchases"], "summary": "更新采购单（仅草稿/已拒绝）", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}, {"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/UpdatePurchaseRequest"}}], "responses": {"200": {"description": "OK"}}}, "delete": {"tags": ["purchases"], "summary": "删除采购单（仅草稿）", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/purchases/{id}/submit": {"put": {"tags": ["purchases"], "summary": "提交审批", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/purchases/{id}/approve": {"put": {"tags": ["purchases"], "summary": "审批通过", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/purchases/{id}/reject": {"put": {"tags": ["purchases"], "summary": "审批拒绝", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/purchases/{id}/complete": {"put": {"tags": ["purchases"], "summary": "标记完成（更新库存）", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}},
    "/stats/cost/by-category": {"get": {"tags": ["stats"], "summary": "按分类统计采购成本", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "start_date", "type": "string"}, {"in": "query", "name": "end_date", "type": "string"}], "responses": {"200": {"description": "OK"}}}},
    "/stats/cost/by-supplier": {"get": {"tags": ["stats"], "summary": "按供应商统计采购总额", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "start_date", "type": "string"}, {"in": "query", "name": "end_date", "type": "string"}], "responses": {"200": {"description": "OK"}}}},
    "/stats/cost/trend": {"get": {"tags": ["stats"], "summary": "成本趋势", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "period", "type": "string"}, {"in": "query", "name": "start_date", "type": "string"}, {"in": "query", "name": "end_date", "type": "string"}], "responses": {"200": {"description": "OK"}}}},
    "/stats/cost/summary": {"get": {"tags": ["stats"], "summary": "成本概览", "security": [{"BearerAuth": []}], "responses": {"200": {"description": "OK"}}}},
    "/logs": {"get": {"tags": ["logs"], "summary": "操作日志列表", "security": [{"BearerAuth": []}], "parameters": [{"in": "query", "name": "page", "type": "integer"}, {"in": "query", "name": "page_size", "type": "integer"}, {"in": "query", "name": "user_id", "type": "integer"}, {"in": "query", "name": "action", "type": "string"}, {"in": "query", "name": "target_type", "type": "string"}, {"in": "query", "name": "start_date", "type": "string"}, {"in": "query", "name": "end_date", "type": "string"}], "responses": {"200": {"description": "OK"}}}},
    "/logs/{id}": {"get": {"tags": ["logs"], "summary": "日志详情", "security": [{"BearerAuth": []}], "parameters": [{"in": "path", "name": "id", "required": true, "type": "integer"}], "responses": {"200": {"description": "OK"}}}}
  },
  "definitions": {
    "LoginRequest": {"type": "object", "required": ["username", "password"], "properties": {"username": {"type": "string"}, "password": {"type": "string"}}},
    "RefreshRequest": {"type": "object", "required": ["refresh_token"], "properties": {"refresh_token": {"type": "string"}}},
    "ChangePasswordRequest": {"type": "object", "required": ["old_password", "new_password"], "properties": {"old_password": {"type": "string"}, "new_password": {"type": "string"}}},
    "CreateUserRequest": {"type": "object", "required": ["username", "password", "role"], "properties": {"username": {"type": "string"}, "password": {"type": "string"}, "role": {"type": "string"}}},
    "UpdateUserRequest": {"type": "object", "required": ["role"], "properties": {"password": {"type": "string"}, "role": {"type": "string"}}},
    "CreateSupplierRequest": {"type": "object", "required": ["name", "contact_person", "phone", "address", "categories"], "properties": {"name": {"type": "string"}, "contact_person": {"type": "string"}, "phone": {"type": "string"}, "email": {"type": "string"}, "address": {"type": "string"}, "categories": {"type": "array", "items": {"type": "string"}}, "rating": {"type": "number"}, "status": {"type": "string"}}},
    "UpdateSupplierRequest": {"type": "object", "required": ["name", "contact_person", "phone", "address", "categories"], "properties": {"name": {"type": "string"}, "contact_person": {"type": "string"}, "phone": {"type": "string"}, "email": {"type": "string"}, "address": {"type": "string"}, "categories": {"type": "array", "items": {"type": "string"}}, "rating": {"type": "number"}}},
    "SupplierStatusRequest": {"type": "object", "required": ["status"], "properties": {"status": {"type": "string"}}},
    "CreateInventoryRequest": {"type": "object", "required": ["name", "category", "supplier_id", "batch_no", "quantity", "unit", "expiry_date", "storage_location"], "properties": {"name": {"type": "string"}, "category": {"type": "string"}, "supplier_id": {"type": "integer"}, "batch_no": {"type": "string"}, "quantity": {"type": "number"}, "unit": {"type": "string"}, "min_threshold": {"type": "number"}, "expiry_date": {"type": "string"}, "storage_location": {"type": "string"}}},
    "UpdateInventoryRequest": {"type": "object", "required": ["name", "category", "supplier_id", "unit", "expiry_date", "storage_location"], "properties": {"name": {"type": "string"}, "category": {"type": "string"}, "supplier_id": {"type": "integer"}, "unit": {"type": "string"}, "min_threshold": {"type": "number"}, "expiry_date": {"type": "string"}, "storage_location": {"type": "string"}}},
    "QuantityAdjustRequest": {"type": "object", "required": ["delta"], "properties": {"delta": {"type": "number"}, "reason": {"type": "string"}}},
    "PurchaseItemRequest": {"type": "object", "required": ["inventory_item_id", "quantity", "unit_price"], "properties": {"inventory_item_id": {"type": "integer"}, "quantity": {"type": "number"}, "unit_price": {"type": "number"}}},
    "CreatePurchaseRequest": {"type": "object", "required": ["supplier_id", "items"], "properties": {"supplier_id": {"type": "integer"}, "notes": {"type": "string"}, "items": {"type": "array", "items": {"$ref": "#/definitions/PurchaseItemRequest"}}}},
    "UpdatePurchaseRequest": {"type": "object", "required": ["supplier_id", "items"], "properties": {"supplier_id": {"type": "integer"}, "notes": {"type": "string"}, "items": {"type": "array", "items": {"$ref": "#/definitions/PurchaseItemRequest"}}}}
  }
}`

// SwaggerInfo 供 gin-swagger 与 /openapi.json 使用。
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "餐饮供应链 API",
	Description:      "餐饮供应链后端管理系统：供应商、库存、采购、成本统计与操作日志。",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
