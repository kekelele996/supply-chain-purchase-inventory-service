# Bug 复现说明

## Bug 是什么

采购单详情与库存/采购单列表的关联数据缺失：明细查不到食材信息、库存列表没有供应商、采购单列表没有创建人。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestGet_ItemsHaveInventory|TestInventoryList_HasSupplier|TestPurchaseList_HasCreator'
```

## 错误信息

```
--- FAIL: TestGet_ItemsHaveInventory
    items should have inventory item loaded: ...
--- FAIL: TestInventoryList_HasSupplier
    inventory should have supplier loaded: ...
--- FAIL: TestPurchaseList_HasCreator
    orders should have creator loaded: ...
```
