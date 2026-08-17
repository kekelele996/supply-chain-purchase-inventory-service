# Bug 复现说明

## Bug 是什么

采购单金额计算错误：明细小计与总金额精度不对（四舍五入变成截断、小计先对数量取整再乘单价），列表按 id 升序排列（最新单在后），数量文案丢失小数位，且 rejected 被排除出合法采购单状态。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestPurchaseCreate_AmountPrecision|TestPurchaseList_NewestFirst|TestQuantify_Decimals|TestIsValidRejectedStatus' -count=1
```

## 错误信息

```
--- FAIL: TestPurchaseCreate_AmountPrecision
    total = 3.14, want 3.15
--- FAIL: TestPurchaseList_NewestFirst
    first id = 1, want 2
--- FAIL: TestQuantify_Decimals
    Quantify = "6 kg", want 5.5 kg
```
