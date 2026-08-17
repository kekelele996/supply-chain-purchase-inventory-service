package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

// ---- fake UserRepository ----

type fakeUserRepo struct {
	mu    sync.Mutex
	byID  map[uint]*model.User
	next  uint
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byID: map[uint]*model.User{}, next: 1}
}

func (f *fakeUserRepo) Create(ctx context.Context, u *model.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range f.byID {
		if v.Username == u.Username {
			return repository.ErrDuplicateKey
		}
	}
	u.ID = f.next
	f.next++
	f.byID[u.ID] = u
	return nil
}

func (f *fakeUserRepo) Update(ctx context.Context, u *model.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byID[u.ID] = u
	return nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id uint) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.byID, id)
	return nil
}

func (f *fakeUserRepo) FindByID(ctx context.Context, id uint) (*model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.byID {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeUserRepo) List(ctx context.Context, page, pageSize int, username, role string) ([]model.User, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []model.User
	for _, u := range f.byID {
		out = append(out, *u)
	}
	return out, int64(len(out)), nil
}

// ---- fake TokenStore ----

type fakeTokenStore struct {
	mu     sync.Mutex
	tokens map[string]uint
	ttl    map[string]time.Time
}

func newFakeTokenStore() *fakeTokenStore {
	return &fakeTokenStore{tokens: map[string]uint{}, ttl: map[string]time.Time{}}
}

func (f *fakeTokenStore) SetRefreshToken(ctx context.Context, token string, userID uint, ttl time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tokens[token] = userID
	f.ttl[token] = time.Now().Add(ttl)
	return nil
}

func (f *fakeTokenStore) GetRefreshTokenUserID(ctx context.Context, token string) (uint, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if exp, ok := f.ttl[token]; ok && time.Now().After(exp) {
		delete(f.tokens, token)
		return 0, nil
	}
	return f.tokens[token], nil
}

func (f *fakeTokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.tokens, token)
	delete(f.ttl, token)
	return nil
}

// ---- fake SupplierRepository ----

type fakeSupplierRepo struct {
	mu    sync.Mutex
	byID  map[uint]*model.Supplier
	next  uint
}

func newFakeSupplierRepo() *fakeSupplierRepo {
	return &fakeSupplierRepo{byID: map[uint]*model.Supplier{}, next: 1}
}

func (f *fakeSupplierRepo) Create(ctx context.Context, s *model.Supplier) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range f.byID {
		if v.Name == s.Name {
			return repository.ErrDuplicateKey
		}
	}
	s.ID = f.next
	f.next++
	f.byID[s.ID] = s
	return nil
}

func (f *fakeSupplierRepo) Update(ctx context.Context, s *model.Supplier) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byID[s.ID] = s
	return nil
}

func (f *fakeSupplierRepo) SoftDelete(ctx context.Context, id uint) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.byID, id)
	return nil
}

func (f *fakeSupplierRepo) FindByID(ctx context.Context, id uint) (*model.Supplier, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s, ok := f.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return s, nil
}

func (f *fakeSupplierRepo) FindByName(ctx context.Context, name string) (*model.Supplier, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, s := range f.byID {
		if s.Name == name {
			return s, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeSupplierRepo) List(ctx context.Context, page, pageSize int, status, category, name string) ([]model.Supplier, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []model.Supplier
	for _, s := range f.byID {
		out = append(out, *s)
	}
	return out, int64(len(out)), nil
}

// ---- fake InventoryRepository ----

type fakeInventoryRepo struct {
	mu    sync.Mutex
	byID  map[uint]*model.InventoryItem
	next  uint
}

func newFakeInventoryRepo() *fakeInventoryRepo {
	return &fakeInventoryRepo{byID: map[uint]*model.InventoryItem{}, next: 1}
}

func (f *fakeInventoryRepo) Create(ctx context.Context, item *model.InventoryItem) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range f.byID {
		if v.BatchNo == item.BatchNo {
			return repository.ErrDuplicateKey
		}
	}
	item.ID = f.next
	f.next++
	f.byID[item.ID] = item
	return nil
}

func (f *fakeInventoryRepo) Update(ctx context.Context, item *model.InventoryItem) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byID[item.ID] = item
	return nil
}

func (f *fakeInventoryRepo) Delete(ctx context.Context, id uint) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.byID, id)
	return nil
}

func (f *fakeInventoryRepo) FindByID(ctx context.Context, id uint) (*model.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return item, nil
}

func (f *fakeInventoryRepo) FindByBatchNo(ctx context.Context, batchNo string) (*model.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, item := range f.byID {
		if item.BatchNo == batchNo {
			return item, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeInventoryRepo) List(ctx context.Context, page, pageSize int, status, category string, supplierID uint, name string) ([]model.InventoryItem, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []model.InventoryItem
	for _, item := range f.byID {
		out = append(out, *item)
	}
	return out, int64(len(out)), nil
}

func (f *fakeInventoryRepo) ListBySupplier(ctx context.Context, supplierID uint, page, pageSize int) ([]model.InventoryItem, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []model.InventoryItem
	for _, item := range f.byID {
		if item.SupplierID == supplierID {
			out = append(out, *item)
		}
	}
	return out, int64(len(out)), nil
}

func (f *fakeInventoryRepo) ListAlerts(ctx context.Context) ([]model.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []model.InventoryItem
	for _, item := range f.byID {
		if item.Quantity < item.MinThreshold {
			out = append(out, *item)
		}
	}
	return out, nil
}

func (f *fakeInventoryRepo) AdjustQuantity(ctx context.Context, id uint, delta float64, newStatus constants.InventoryStatus) (*model.InventoryItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	if item.Quantity+delta < 0 {
		return nil, repository.ErrConflict
	}
	item.Quantity += delta
	item.Status = newStatus
	return item, nil
}

func (f *fakeInventoryRepo) ForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.InventoryItem, error) {
	return f.FindByID(ctx, id)
}

// ---- fake PurchaseRepository ----

type fakePurchaseRepo struct {
	mu    sync.Mutex
	byID  map[uint]*model.PurchaseOrder
	next  uint
}

func newFakePurchaseRepo() *fakePurchaseRepo {
	return &fakePurchaseRepo{byID: map[uint]*model.PurchaseOrder{}, next: 1}
}

func (f *fakePurchaseRepo) Create(ctx context.Context, order *model.PurchaseOrder, items []model.PurchaseOrderItem) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range f.byID {
		if v.OrderNo == order.OrderNo {
			return repository.ErrDuplicateKey
		}
	}
	order.ID = f.next
	f.next++
	order.Items = items
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
	}
	f.byID[order.ID] = order
	return nil
}

func (f *fakePurchaseRepo) UpdateWithItems(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder, items []model.PurchaseOrderItem) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	order.Items = items
	f.byID[order.ID] = order
	return nil
}

func (f *fakePurchaseRepo) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byID[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.byID, id)
	return nil
}

func (f *fakePurchaseRepo) FindByID(ctx context.Context, id uint) (*model.PurchaseOrder, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	o, ok := f.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return o, nil
}

func (f *fakePurchaseRepo) FindByOrderNo(ctx context.Context, orderNo string) (*model.PurchaseOrder, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, o := range f.byID {
		if o.OrderNo == orderNo {
			return o, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakePurchaseRepo) List(ctx context.Context, page, pageSize int, status string, supplierID uint, startDate, endDate string) ([]model.PurchaseOrder, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []model.PurchaseOrder
	for _, o := range f.byID {
		out = append(out, *o)
	}
	return out, int64(len(out)), nil
}

func (f *fakePurchaseRepo) UpdateStatus(ctx context.Context, tx *gorm.DB, order *model.PurchaseOrder) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byID[order.ID] = order
	return nil
}

func (f *fakePurchaseRepo) ForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.PurchaseOrder, error) {
	return f.FindByID(ctx, id)
}

func (f *fakePurchaseRepo) CountByPeriod(ctx context.Context, tx *gorm.DB, start, end time.Time) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return int64(len(f.byID)), nil
}

// ---- helpers ----

func mustUser(t *testing.T, username, password, role string) *model.User {
	t.Helper()
	hash, err := util.HashPassword(password)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	return &model.User{Username: username, PasswordHash: hash, Role: constants.UserRole(role)}
}

var _ repository.UserRepository = (*fakeUserRepo)(nil)
var _ repository.SupplierRepository = (*fakeSupplierRepo)(nil)
var _ repository.InventoryRepository = (*fakeInventoryRepo)(nil)
var _ repository.PurchaseRepository = (*fakePurchaseRepo)(nil)
var _ TokenStore = (*fakeTokenStore)(nil)
