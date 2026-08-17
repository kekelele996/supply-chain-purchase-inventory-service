package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

// UserService 用户管理服务（admin 专用）。
type UserService struct {
	users repository.UserRepository
}

// NewUserService 构造用户管理服务。
func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
}

// Create 创建用户。
func (s *UserService) Create(ctx context.Context, username, password, role string) (*model.User, error) {
	if _, err := s.users.FindByUsername(ctx, username); err == nil {
		return nil, util.NewAppError(constants.CodeUserNameExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextUserExists, username), nil)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("create user: %w", err)
	}
	hash, err := util.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	u := &model.User{Username: username, PasswordHash: hash, Role: constants.NormalizeRole(role)}
	if err := s.users.Create(ctx, u); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, util.NewAppError(constants.CodeUserNameExists, http.StatusConflict, fmt.Sprintf(constants.ErrTextUserExists, username), err)
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

// Update 更新用户（角色与可选密码）。
func (s *UserService) Update(ctx context.Context, id uint, password, role string) (*model.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "用户", id), err)
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	if password != "" {
		hash, herr := util.HashPassword(password)
		if herr != nil {
			return nil, fmt.Errorf("update user: %w", herr)
		}
		u.PasswordHash = hash
	}
	u.Role = constants.NormalizeRole(role)
	if err := s.users.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

// Delete 删除用户。
func (s *UserService) Delete(ctx context.Context, id uint) error {
	if err := s.users.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeUserNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "用户", id), err)
		}
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// Get 获取用户详情。
func (s *UserService) Get(ctx context.Context, id uint) (*model.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, http.StatusNotFound, fmt.Sprintf(constants.ErrTextNotFound, "用户", id), err)
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// List 分页查询用户。
func (s *UserService) List(ctx context.Context, page, pageSize int, username, role string) ([]model.User, int64, error) {
	list, total, err := s.users.List(ctx, page, pageSize, username, role)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return list, total, nil
}
