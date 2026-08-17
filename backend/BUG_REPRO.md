# Bug 复现说明

## Bug 是什么

库存状态计算错误：余量等于预警阈值被标为缺货、还有一天才过期被算成过期、完成采购单后库存状态未重算，库存状态中文文案也错。

## 如何触发

```bash
cd backend
go test ./... -run 'TestComputeStatus_AtThreshold|TestComputeStatus_TomorrowNotExpired|TestComplete_UpdatesInventoryStatus|TestInventoryStatusText_Normal' -count=1
```

## 错误信息

```
--- FAIL: TestComputeStatus_AtThreshold
    at threshold = low, want normal
--- FAIL: TestComputeStatus_TomorrowNotExpired
    tomorrow expiry = expired, want normal
--- FAIL: TestComplete_UpdatesInventoryStatus
    status = low, want normal
```
