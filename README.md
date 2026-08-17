# 餐饮供应链 API（wje-108 / supplychain）

面向餐饮企业的供应链后端管理系统：提供供应商管理、食材库存管理、采购订单审批流转、成本核算与操作审计日志等核心能力，通过 RESTful API 为前端或第三方系统提供数据服务。

## Docker 部署（首选）

```bash
# 1. 进入项目根目录（任意目录名，包括中文目录名均可）
cd wje-108

# 2. 一键启动全部服务（MySQL + Redis + MinIO + 后端）
docker compose up -d --build

# 3. 查看状态（等待全部容器 healthy）
docker compose ps

# 4. 验证健康检查
curl -sS http://localhost:29608/healthz

# 5. 停止并清理（删除数据卷）
docker compose down -v --remove-orphans
```

启动后访问：

| 服务 | 地址 |
| --- | --- |
| 后端 API | http://localhost:29608/api/v1 |
| Swagger UI | http://localhost:29608/docs （备用 http://localhost:29608/swagger/index.html） |
| OpenAPI JSON | http://localhost:29608/openapi.json |
| 健康检查 | http://localhost:29608/healthz |
| MinIO 控制台 | http://localhost:47023 （账号 minioadmin / minioadmin123） |

## 本地开发（备选）

```bash
cd backend
go mod tidy
go run ./cmd/server        # 默认监听 :8080，可通过 SERVER_PORT 覆盖
go build ./...
go vet ./...
go test ./...
```

本地运行需要自行准备 MySQL（`DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME`）与 Redis（`REDIS_ADDR`），见 `.env.example`。

## 项目主要功能

- **认证授权**：JWT 登录/刷新 + API Key 双认证，RBAC 角色权限（operator / manager / admin）。
- **供应商管理**：增删改查、软删除、状态流转（active / suspended / blacklisted）、按供应商查看库存与采购单。
- **库存管理**：入库/出库/调整余量、批次号唯一、保质期与预警阈值自动计算状态（normal / low / expired）。
- **采购管理**：采购单完整状态机（draft → pending_approval → approved/rejected → completed），审批权限控制（不可自审），完成时事务内联动更新库存。
- **成本统计**：按食材分类 / 按供应商统计采购成本、成本趋势、成本概览。
- **操作审计**：所有写操作自动记录操作日志，支持按用户/动作/对象/时间范围筛选。
- **横切能力**：请求 ID 追踪、Redis 限流、统一错误处理、结构化日志（slog）。

## 技术栈

| 层 | 技术栈 |
| --- | --- |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7 |
| 对象存储 | MinIO（预留） |
| 认证 | JWT + API Key + RBAC |
| 日志 | log/slog 结构化日志 |
| 参数校验 | github.com/go-playground/validator/v10 |
| 接口文档 | Swagger/OpenAPI 2.0 |

## 项目目录结构

```
wje-108/
├── backend/                       # Go 后端
│   ├── cmd/server/main.go         # 入口：装配依赖、启动服务
│   ├── internal/
│   │   ├── config/config.go       # 环境变量配置解析
│   │   ├── database/              # MySQL 连接、AutoMigrate、种子数据
│   │   ├── model/                 # 每个实体一个文件（user/supplier/inventory/purchase_order/...）
│   │   ├── dto/                   # 每个实体一个 DTO 文件
│   │   ├── repository/            # 每个实体一个仓储文件（含哨兵错误、行锁）
│   │   ├── service/               # 每个实体一个服务文件（状态机、事务）
│   │   ├── handler/               # 每个实体一个处理器文件
│   │   ├── router/                # 每个实体一个路由注册文件
│   │   ├── middleware/            # auth/rbac/api_key/request_id/audit/rate_limiter/error_handler/operation_logger
│   │   ├── constants/             # 枚举、错误码、文案、日志模板
│   │   └── util/                  # logger/jwt/password/formatters/response/validator 等
│   ├── pkg/strutil/               # 无业务依赖的可复用工具（单号生成等）
│   ├── docs/                      # OpenAPI 2.0 规范（Swagger UI 数据）
│   ├── migrations/                # SQL 迁移脚本（与 GORM 模型一致）
│   ├── Dockerfile                 # Go 多阶段构建
│   ├── go.mod / go.sum
├── database/init.sql              # MySQL 首次启动初始化脚本（建表 + 种子数据）
├── docker-compose.yml
├── .env / .env.example
├── README.md
└── output/execution.md            # 执行验证报告
```

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `supplychain` | Compose 项目名/容器名前缀 |
| `BACKEND_PORT` | `8080` | 后端宿主机端口（本项目分配 `29608`） |
| `DB_PORT` | `3306` | MySQL 宿主机端口（本项目分配 `57603`） |
| `REDIS_PORT` | `6379` | Redis 宿主机端口（本项目分配 `46320`） |
| `MINIO_PORT` / `MINIO_CONSOLE_PORT` | `9000`/`9001` | MinIO API/控制台端口（本项目 `47022`/`47023`） |
| `DB_NAME` / `DB_USER` / `DB_PASSWORD` | `supplychain_db` / `supplychain_user` / `supplychain_pwd` | MySQL 库名/账号/密码 |
| `REDIS_ADDR` / `REDIS_PASSWORD` | `redis:6379` / 空 | Redis 连接地址与密码 |
| `JWT_SECRET` | 占位串 | JWT 签名密钥（生产务必更换） |
| `API_KEY_SECRET` | 占位串 | `X-API-Key` 校验密钥（生产务必更换） |
| `RATE_LIMIT_PER_MINUTE` | `600` | 单 IP 每分钟请求上限 |

## API 清单（统一前缀 /api/v1，响应 {code, message, data}）

### 认证 /auth
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| POST | /auth/login | 登录，返回 access_token + refresh_token | 公开 |
| POST | /auth/refresh | 刷新 token（轮换） | 登录 |
| GET | /auth/me | 当前用户信息 | 登录 |
| PUT | /auth/password | 修改密码 | 登录 |

### 用户管理 /users（admin）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /users | 用户列表（分页 + username/role 筛选） |
| POST | /users | 创建用户 |
| GET | /users/{id} | 用户详情 |
| PUT | /users/{id} | 更新用户 |
| DELETE | /users/{id} | 删除用户 |

### 供应商 /suppliers
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /suppliers | 列表（分页 + status/category/name 筛选） | 登录 |
| POST | /suppliers | 新增供应商 | operator+ |
| GET | /suppliers/{id} | 详情 | 登录 |
| PUT | /suppliers/{id} | 更新 | operator+ |
| DELETE | /suppliers/{id} | 删除（软删除） | admin |
| PUT | /suppliers/{id}/status | 变更状态 | manager+ |
| GET | /suppliers/{id}/inventory | 该供应商库存列表 | 登录 |
| GET | /suppliers/{id}/orders | 该供应商采购单列表 | 登录 |

### 库存 /inventory
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /inventory | 列表（分页 + status/category/supplier_id/name 筛选） | 登录 |
| POST | /inventory | 入库 | operator+ |
| GET | /inventory/alerts | 库存预警（低于阈值/已过期） | 登录 |
| GET | /inventory/{id} | 详情 | 登录 |
| PUT | /inventory/{id} | 更新信息 | operator+ |
| PUT | /inventory/{id}/quantity | 调整余量（delta 正入库负出库） | operator+ |
| DELETE | /inventory/{id} | 删除 | admin |

### 采购 /purchases
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /purchases | 列表（分页 + status/supplier_id/start_date/end_date 筛选） | 登录 |
| POST | /purchases | 创建采购单（含明细） | operator+ |
| GET | /purchases/{id} | 详情（含明细） | 登录 |
| PUT | /purchases/{id} | 更新（仅 draft/rejected） | operator+ |
| DELETE | /purchases/{id} | 删除（仅 draft） | operator+ |
| PUT | /purchases/{id}/submit | 提交审批（draft/rejected → pending_approval） | operator+ |
| PUT | /purchases/{id}/approve | 审批通过（pending_approval → approved） | manager+，不可自审 |
| PUT | /purchases/{id}/reject | 审批拒绝（pending_approval → rejected） | manager+，不可自审 |
| PUT | /purchases/{id}/complete | 完成（approved → completed，事务内更新库存） | operator+ |

### 成本统计 /stats
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /stats/cost/by-category | 按食材分类统计采购成本 |
| GET | /stats/cost/by-supplier | 按供应商统计采购总额 |
| GET | /stats/cost/trend | 成本趋势（period=day/week/month） |
| GET | /stats/cost/summary | 成本概览（总采购额/平均单笔/本月采购额） |

### 操作日志 /logs（manager+）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /logs | 日志列表（分页 + user_id/action/target_type/date_range 筛选） |
| GET | /logs/{id} | 日志详情 |

## curl 调用示例

```bash
BASE=http://localhost:29608/api/v1

# 1) 登录获取 token（默认种子账号：admin/admin123、manager/manager123、operator/operator123）
TOKEN=$(curl -sS -X POST $BASE/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.access_token')
echo "TOKEN=$TOKEN"

# 2) 携带 JWT 获取当前用户
curl -sS $BASE/auth/me -H "Authorization: Bearer $TOKEN" | jq

# 3) 供应商列表
curl -sS "$BASE/suppliers?page=1&page_size=5" -H "Authorization: Bearer $TOKEN" | jq

# 4) 新增供应商
curl -sS -X POST $BASE/suppliers -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"示例供应商","contact_person":"张三","phone":"13800009999","email":"z@example.com","address":"某市某路1号","categories":["蔬菜","水果"],"rating":4.0}' | jq

# 5) 库存列表与预警
curl -sS "$BASE/inventory?status=low" -H "Authorization: Bearer $TOKEN" | jq
curl -sS $BASE/inventory/alerts -H "Authorization: Bearer $TOKEN" | jq

# 6) 创建采购单（items 中 inventory_item_id 为库存 ID）
curl -sS -X POST $BASE/purchases -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"supplier_id":1,"notes":"首批采购","items":[{"inventory_item_id":1,"quantity":20,"unit_price":5.5}]}' | jq

# 7) 提交/审批/完成（审批需 manager/admin，且不能审批自己的单）
curl -sS -X PUT $BASE/purchases/1/submit -H "Authorization: Bearer $TOKEN" | jq
MTOKEN=$(curl -sS -X POST $BASE/auth/login -H 'Content-Type: application/json' -d '{"username":"manager","password":"manager123"}' | jq -r '.data.access_token')
curl -sS -X PUT $BASE/purchases/1/approve -H "Authorization: Bearer $MTOKEN" | jq
curl -sS -X PUT $BASE/purchases/1/complete -H "Authorization: Bearer $TOKEN" | jq

# 8) 成本统计
curl -sS $BASE/stats/cost/summary -H "Authorization: Bearer $TOKEN" | jq
curl -sS "$BASE/stats/cost/by-category" -H "Authorization: Bearer $TOKEN" | jq

# 9) 操作日志（manager+）
curl -sS "$BASE/logs?page=1&page_size=10" -H "Authorization: Bearer $MTOKEN" | jq

# 10) API Key 双认证示例
curl -sS $BASE/suppliers -H "X-API-Key: change_me_to_a_long_random_api_key_supplychain_2026" | jq
```

## Docker 部署说明

- **端口映射**：后端 `29608 → 8080`，MySQL `57603 → 3306`，Redis `46320 → 6379`，MinIO `47022 → 9000`、`47023 → 9001`（均可通过 `.env` 修改）。
- **数据持久化**：MySQL/Redis/MinIO 使用命名卷 `db_data` / `redis_data` / `minio_data`，`docker compose down -v` 才会删除数据。
- **初始化**：MySQL 容器首次启动自动执行 `database/init.sql`（建库建表 + 种子用户/供应商/库存）；后端启动时执行 GORM AutoMigrate 增量同步并在用户表为空时兜底播种。
- **依赖顺序**：后端 `depends_on` MySQL/Redis 健康检查通过后才启动；后端自身提供 `/healthz` 健康检查。
- **常见问题**：
  - 端口被占用：修改 `.env` 中的端口后重新 `docker compose up -d`。
  - 容器未 healthy：`docker compose logs backend` / `docker compose logs db` 查看日志。
  - 修改了代码：`docker compose up -d --build backend` 重新构建。

## 枚举出现位置清单

### 1) PurchaseOrderStatus 采购单状态（draft / pending_approval / approved / rejected / completed）
- 定义：`backend/internal/constants/enums.go`
- 模型：`backend/internal/model/purchase_order.go`
- DTO：`backend/internal/dto/purchase_dto.go`
- 状态机：`backend/internal/service/purchase_service.go`（Create/Update/Submit/Approve/Reject/Complete/Delete 流转校验）
- 校验：`backend/internal/constants/enums.go`（IsValidPurchaseOrderStatus）、`backend/internal/util/formatters.go`（OrderStatusText）
- 错误码：`backend/internal/constants/error_codes.go`（CodeOrderStateInvalid / CodeOrderCannotSelfApprove / CodeOrderSupplierSuspended）
- 日志模板：`backend/internal/constants/log_templates.go`（LogTplOrderSubmit/Approve/Reject/Complete/StateInvalid）
- 文案：`backend/internal/constants/messages.go`（MsgOrderSubmitted/Approved/Rejected/Completed）
- 统计口径：`backend/internal/service/stats_service.go`（approved/completed 才计入成本）
- 路由权限：`backend/internal/router/purchases.go`
- 迁移/初始化：`backend/migrations/0001_init.sql`、`database/init.sql`

### 2) SupplierStatus 供应商状态（active / suspended / blacklisted）
- 定义：`backend/internal/constants/enums.go`
- 模型：`backend/internal/model/supplier.go`
- DTO：`backend/internal/dto/supplier_dto.go`（SupplierStatusRequest oneof 校验）
- 业务校验：`backend/internal/service/supplier_service.go`（ChangeStatus）、`backend/internal/service/purchase_service.go`（非 active 禁止创建采购单）
- 文案：`backend/internal/util/formatters.go`（SupplierStatusText）、`backend/internal/constants/messages.go`
- 错误码：`backend/internal/constants/error_codes.go`（CodeSupplierForbidden / CodeOrderSupplierSuspended）
- 日志模板：`backend/internal/constants/log_templates.go`（LogTplSupplierStatus）

### 3) InventoryStatus 库存状态（normal / low / expired）
- 定义：`backend/internal/constants/enums.go`
- 模型：`backend/internal/model/inventory_item.go`
- 计算：`backend/internal/service/inventory_service.go`（computeInventoryStatus，入库/调整余量/完成采购联动时重算）
- 仓储：`backend/internal/repository/inventory_repository.go`（ListAlerts 按余量/保质期查询）
- 文案：`backend/internal/util/formatters.go`（InventoryStatusText）
- 日志模板：`backend/internal/constants/log_templates.go`（LogTplInventoryAdjust）

### 4) InventoryCategory 食材分类（vegetable / meat / seafood / seasoning / staple / dry_goods / beverage）
- 定义：`backend/internal/constants/enums.go`
- 模型：`backend/internal/model/inventory_item.go`
- DTO：`backend/internal/dto/inventory_dto.go`（oneof 校验）
- 成本统计分组：`backend/internal/service/stats_service.go`（ByCategory 按 category 分组）
- 文案：`backend/internal/util/formatters.go`（InventoryCategoryText）

### 5) UserRole 用户角色（operator / manager / admin）
- 定义：`backend/internal/constants/enums.go`
- 模型：`backend/internal/model/user.go`
- DTO：`backend/internal/dto/user_dto.go`（oneof 校验）
- RBAC：`backend/internal/middleware/rbac.go`、`backend/internal/router/*.go` 各路由角色声明
- 业务校验：`backend/internal/service/purchase_service.go`（审批角色 + 不可自审）
- 文案：`backend/internal/util/formatters.go`（RoleText）

### 6) OperationAction 操作动作（create / update / delete / approve / reject / login / logout / submit / complete）
- 定义：`backend/internal/constants/enums.go`
- 模型：`backend/internal/model/operation_log.go`
- 中间件：`backend/internal/middleware/operation_logger.go`（写操作自动映射动作）
- 服务：`backend/internal/service/operation_log_service.go`
- 日志模板：`backend/internal/constants/log_templates.go`（LogTplOperationRecorded）

## 横切关注点（触达文件层）

1. **API Key + JWT 双认证 + RBAC 权限**
   - 触达：`internal/middleware/auth.go`（双认证入口）、`internal/middleware/api_key.go`（API Key 常量时间比较）、`internal/middleware/rbac.go`（角色校验）、`internal/util/jwt.go`（JWT 签发/解析）、`internal/router/*.go`（路由权限声明）、`internal/service/auth_service.go`（登录/刷新）、`internal/config/config.go`（密钥配置）。
2. **操作审计日志与请求追踪**
   - 触达：`internal/middleware/request_id.go`（X-Request-ID 生成/透传）、`internal/middleware/audit.go`（请求级审计日志）、`internal/middleware/operation_logger.go`（写操作异步落库）、`internal/service/operation_log_service.go`、`internal/repository/operation_log_repository.go`、`internal/model/operation_log.go`、`internal/constants/log_templates.go`。
3. **接口限流**
   - 触达：`internal/middleware/rate_limiter.go`（Redis 固定窗口）、`internal/config/config.go`（RATE_LIMIT_PER_MINUTE）、`internal/constants/error_codes.go`（CodeRateLimited）、`internal/constants/log_templates.go`（LogTplRateLimited）。

## 后端中间件清单

| 中间件 | 文件 | 说明 |
| --- | --- | --- |
| 请求 ID | `internal/middleware/request_id.go` | 生成/透传 X-Request-ID |
| 统一异常处理 | `internal/middleware/error_handler.go` | panic 恢复 + 统一错误响应 |
| 请求审计 | `internal/middleware/audit.go` | 记录 method/path/status/latency/user |
| 限流 | `internal/middleware/rate_limiter.go` | Redis 固定窗口限流 |
| 认证 | `internal/middleware/auth.go` | JWT + API Key 双认证 |
| RBAC | `internal/middleware/rbac.go` | 角色权限校验 |
| 操作日志 | `internal/middleware/operation_logger.go` | 写操作异步审计 |

## 文件结构强制清单

- 每个核心实体独立贯穿：`model → dto → repository → service → handler → router → constants`，文件一一对应，**严禁合并职责到单一文件**（例如不允许把多个实体的 model/repository/service/handler 写进同一文件，也不允许 main.go 直接写业务逻辑）。
- 依赖方向严格单向：`handler → service → repository → model`，构造器注入，禁止包循环。
- 多步写操作放入 service 事务，并发场景使用 `SELECT ... FOR UPDATE`（`internal/repository/clause.go`）。

## 屎山代码设计要求（牵一发动全身）

为验证跨文件协同改动能力，本工程按提示词要求保留了以下“合理的屎山”耦合点：

1. **日志模块单独管理但全栈引用**：`internal/util/logger.go` 统一封装 slog；所有 handler/service/middleware 引用 logger；≥25 条日志模板集中在 `internal/constants/log_templates.go`；业务字段/状态变更需同步模板与调用处。
2. **异常信息分散且层层透传**：错误码集中在 `internal/constants/error_codes.go`，service/handler 手动拼接 message，handler 再次包装；message 中包含实体名、字段名、角色名（见 `internal/constants/messages.go` 的 ErrText* 模板）。
3. **常量/工具类多处耦合**：`internal/util/formatters.go` 同时包含日期、状态文本、类型文本格式化；`internal/constants/messages.go` 同时包含接口文案、日志文案、错误提示文案。
4. **状态机跨多处定义**：采购单状态流转规则同时存在于 `service/purchase_service.go`、`util/formatters.go`、`constants/log_templates.go`、`constants/error_codes.go`、`constants/messages.go`、`dto/purchase_dto.go`、`router/purchases.go` 等 ≥10 处；新增一个状态需同步修改。
5. **枚举多处重复定义且被多处引用**：核心枚举在 constants、DTO、模型、日志模板、错误码、formatters 中同时存在（见上文枚举出现位置清单）；新增一个枚举值需修改 ≥10 个文件。

## License

MIT License
