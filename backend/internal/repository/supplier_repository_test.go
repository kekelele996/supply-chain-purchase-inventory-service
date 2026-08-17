package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

func newSupplier(name string) *model.Supplier {
	return &model.Supplier{
		Name:          name,
		ContactPerson: "张三",
		Phone:         "13800000000",
		Email:         "a@b.cn",
		Address:       "测试地址",
		Categories:    model.StringList{"蔬菜", "肉类"},
		Rating:        4.0,
		Status:        constants.SupplierActive,
	}
}

func TestSupplierRepository_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewSupplierRepository(db)
	ctx := context.Background()
	s := newSupplier("供应商A")
	if err := repo.Create(ctx, s); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if got.Name != "供应商A" || len(got.Categories) != 2 {
		t.Fatalf("unexpected supplier: %+v", got)
	}
}

func TestSupplierRepository_ListWithCategoryFilter(t *testing.T) {
	db := newTestDB(t)
	repo := NewSupplierRepository(db)
	ctx := context.Background()
	a := newSupplier("蔬菜供应商")
	b := newSupplier("肉类供应商")
	b.Categories = model.StringList{"肉类"}
	if err := repo.Create(ctx, a); err != nil {
		t.Fatalf("create a: %v", err)
	}
	if err := repo.Create(ctx, b); err != nil {
		t.Fatalf("create b: %v", err)
	}
	list, total, err := repo.List(ctx, 1, 10, "", "肉类", "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 2 || len(list) != 2 {
		t.Fatalf("expected 2 matches for 肉类, got total=%d len=%d", total, len(list))
	}
	list, total, err = repo.List(ctx, 1, 10, "", "海鲜", "")
	if err != nil {
		t.Fatalf("list seafood: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected 0 matches for 海鲜, got %d", total)
	}
}

func TestSupplierRepository_SoftDelete(t *testing.T) {
	db := newTestDB(t)
	repo := NewSupplierRepository(db)
	ctx := context.Background()
	s := newSupplier("待删除")
	if err := repo.Create(ctx, s); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.SoftDelete(ctx, s.ID); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, s.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found after soft delete, got %v", err)
	}
}
