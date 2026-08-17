# Bug 复现说明

## Bug 是什么

采购单审批状态机错乱：创建人能审批自己的单子、草稿单能直接被拒绝、被拒绝的单子不能重新提交，审批人关联数据查询不到，状态文案也错。

## 如何触发

```bash
cd backend
go test -run 'TestApprove_SelfApproveBlocked|TestReject_DraftBlocked|TestSubmit_RejectedResubmitAllowed|TestOrderStatusText_Approved|TestIsValidCompletedStatus|TestGet_HasApprover' ./... -count=1
```

## 错误信息

```
--- FAIL: TestApprove_SelfApproveBlocked
    expected self-approve error
--- FAIL: TestReject_DraftBlocked
    expected state invalid for draft reject
--- FAIL: TestSubmit_RejectedResubmitAllowed
    resubmit rejected order: ...
--- FAIL: TestGet_HasApprover
    approver should be loaded
```
