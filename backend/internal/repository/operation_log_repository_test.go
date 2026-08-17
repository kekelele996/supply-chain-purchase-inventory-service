package repository

import (
	"context"
	"testing"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

func TestOperationLogRepository_CreateAndList(t *testing.T) {
	db := newTestDB(t)
	repo := NewOperationLogRepository(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if err := repo.Create(ctx, &model.OperationLog{
			UserID:     1,
			Action:     constants.ActionCreate,
			TargetType: "supplier",
			TargetID:   uint(i + 1),
			Detail:     `{"method":"POST"}`,
		}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	list, total, err := repo.List(ctx, 1, 10, 1, "create", "supplier", "", "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 || len(list) != 3 {
		t.Fatalf("expected 3, got total=%d len=%d", total, len(list))
	}
	entry, err := repo.FindByID(ctx, list[0].ID)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if entry.TargetType != "supplier" {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}
