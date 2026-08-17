package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/util"
)

// Seed 首次启动时写入默认用户与示例业务数据（幂等：已存在则跳过）。
func Seed(ctx context.Context, db *gorm.DB) error {
	seeded, err := seedUsers(ctx, db)
	if err != nil {
		return err
	}
	supplierCount, err := seedSuppliersAndInventory(ctx, db)
	if err != nil {
		return err
	}
	util.Logger.Info("database seed finished", "users", seeded, "suppliers", supplierCount)
	return nil
}

func seedUsers(ctx context.Context, db *gorm.DB) (int, error) {
	users := []struct {
		Username string
		Password string
		Role     constants.UserRole
	}{
		{Username: "admin", Password: "admin123", Role: constants.RoleAdmin},
		{Username: "manager", Password: "manager123", Role: constants.RoleManager},
		{Username: "operator", Password: "operator123", Role: constants.RoleOperator},
	}
	count := 0
	for _, u := range users {
		var existing int64
		if err := db.WithContext(ctx).Model(&model.User{}).Where("username = ?", u.Username).Count(&existing).Error; err != nil {
			return count, fmt.Errorf("seed users: %w", err)
		}
		if existing > 0 {
			continue
		}
		hash, err := util.HashPassword(u.Password)
		if err != nil {
			return count, fmt.Errorf("seed users: %w", err)
		}
		if err := db.WithContext(ctx).Create(&model.User{Username: u.Username, PasswordHash: hash, Role: u.Role}).Error; err != nil {
			return count, fmt.Errorf("seed users: %w", err)
		}
		count++
	}
	return count, nil
}

func seedSuppliersAndInventory(ctx context.Context, db *gorm.DB) (int, error) {
	var supCount int64
	if err := db.WithContext(ctx).Model(&model.Supplier{}).Count(&supCount).Error; err != nil {
		return 0, fmt.Errorf("seed suppliers: %w", err)
	}
	if supCount > 0 {
		return int(supCount), nil
	}
	now := time.Now()
	suppliers := []model.Supplier{
		{Name: "绿源蔬菜基地", ContactPerson: "王强", Phone: "13800000001", Email: "wang@greenfarm.cn", Address: "北京市朝阳区蔬菜基地1号", Categories: model.StringList{"蔬菜", "水果"}, Rating: 4.5, Status: constants.SupplierActive},
		{Name: "鲜丰肉类供应", ContactPerson: "李娜", Phone: "13800000002", Email: "lina@meat.cn", Address: "上海市浦东新区肉联厂路2号", Categories: model.StringList{"肉类"}, Rating: 4.0, Status: constants.SupplierActive},
		{Name: "海味源水产", ContactPerson: "赵海", Phone: "13800000003", Email: "zhaohai@seafood.cn", Address: "广州市海珠区水产市场3号", Categories: model.StringList{"海鲜"}, Rating: 3.5, Status: constants.SupplierActive},
		{Name: "百味调料行", ContactPerson: "陈香", Phone: "13800000004", Email: "chenxiang@seasoning.cn", Address: "成都市锦江区调料市场4号", Categories: model.StringList{"调料", "干货"}, Rating: 3.0, Status: constants.SupplierSuspended},
	}
	if err := db.WithContext(ctx).Create(&suppliers).Error; err != nil && !isDup(err) {
		return 0, fmt.Errorf("seed suppliers: %w", err)
	}
	_ = now
	inventory := []model.InventoryItem{
		{Name: "有机番茄", Category: constants.CategoryVegetable, SupplierID: suppliers[0].ID, BatchNo: "VEG-20260801", Quantity: 120, Unit: "kg", MinThreshold: 30, ExpiryDate: time.Now().Add(5 * 24 * time.Hour), StorageLocation: "A区-01架", Status: constants.InventoryNormal},
		{Name: "猪五花肉", Category: constants.CategoryMeat, SupplierID: suppliers[1].ID, BatchNo: "MEAT-20260802", Quantity: 8, Unit: "kg", MinThreshold: 20, ExpiryDate: time.Now().Add(3 * 24 * time.Hour), StorageLocation: "B区-冷冻2号", Status: constants.InventoryLow},
		{Name: "东海带鱼", Category: constants.CategorySeafood, SupplierID: suppliers[2].ID, BatchNo: "SEA-20260720", Quantity: 50, Unit: "箱", MinThreshold: 10, ExpiryDate: time.Now().Add(-1 * 24 * time.Hour), StorageLocation: "C区-冷藏3号", Status: constants.InventoryExpired},
		{Name: "郫县豆瓣酱", Category: constants.CategorySeasoning, SupplierID: suppliers[3].ID, BatchNo: "SEA-20260701", Quantity: 200, Unit: "瓶", MinThreshold: 50, ExpiryDate: time.Now().Add(60 * 24 * time.Hour), StorageLocation: "D区-调料架", Status: constants.InventoryNormal},
	}
	if err := db.WithContext(ctx).Create(&inventory).Error; err != nil {
		return int(supCount), fmt.Errorf("seed inventory: %w", err)
	}
	return int(supCount), nil
}

func isDup(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
