# Bug 复现说明

## Bug 是什么

库存预警漏掉已过期食材（只按数量低于阈值判断），成本统计把已完成的采购单排除在 total_amount 之外。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestAlerts_IncludeExpired|TestStatsSummary_IncludeCompleted'
```

## 错误信息

```
--- FAIL: TestAlerts_IncludeExpired
    expired item should be included in alerts
--- FAIL: TestStatsSummary_IncludeCompleted
    total = 12, want 62
```
