package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/supplychain/supplychain-api/internal/model"
)

// UserRepository 用户仓储接口。
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, page, pageSize int, username, role string) ([]model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create user: %w", ErrDuplicateKey)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Save(u).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&model.User{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete user id=%d: %w", id, ErrNotFound)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user id=%d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find user id=%d: %w", id, err)
	}
	return &u, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user %q: %w", username, ErrNotFound)
		}
		return nil, fmt.Errorf("find user %q: %w", username, err)
	}
	return &u, nil
}

func (r *userRepository) List(ctx context.Context, page, pageSize int, username, role string) ([]model.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.User{})
	if username != "" {
		q = q.Where("username LIKE ?", "%"+username+"%")
	}
	if role != "" {
		q = q.Where("role = ?", role)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	var users []model.User
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}

func isDuplicateErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique") || strings.Contains(msg, "1062")
}
