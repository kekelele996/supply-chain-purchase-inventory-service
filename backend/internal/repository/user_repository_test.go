package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
)

func TestUserRepository_CreateAndFindByUsername(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	u := &model.User{Username: "alice", PasswordHash: "hash1", Role: constants.RoleOperator}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("find by username: %v", err)
	}
	if got.ID == 0 || got.Username != "alice" {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestUserRepository_CreateDuplicate(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	u := &model.User{Username: "bob", PasswordHash: "hash", Role: constants.RoleManager}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	dup := &model.User{Username: "bob", PasswordHash: "hash2", Role: constants.RoleAdmin}
	if err := repo.Create(ctx, dup); !errors.Is(err, ErrDuplicateKey) {
		t.Fatalf("expected duplicate key, got %v", err)
	}
}

func TestUserRepository_ListFilter(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	for _, u := range []*model.User{
		{Username: "u1", PasswordHash: "h", Role: constants.RoleOperator},
		{Username: "u2", PasswordHash: "h", Role: constants.RoleManager},
		{Username: "u3", PasswordHash: "h", Role: constants.RoleAdmin},
	} {
		if err := repo.Create(ctx, u); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	list, total, err := repo.List(ctx, 1, 10, "u", "manager")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Username != "u2" {
		t.Fatalf("unexpected list result total=%d list=%+v", total, list)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	u := &model.User{Username: "del", PasswordHash: "h", Role: constants.RoleOperator}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.Delete(ctx, u.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, u.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
