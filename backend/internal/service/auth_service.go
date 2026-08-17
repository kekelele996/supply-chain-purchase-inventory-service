package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/model"
	"github.com/supplychain/supplychain-api/internal/repository"
	"github.com/supplychain/supplychain-api/internal/util"
)

// AuthService 认证服务。
type AuthService struct {
	users      repository.UserRepository
	tokenStore TokenStore
	cfg        AuthServiceConfig
	logger     *slog.Logger
}

// AuthServiceConfig 认证服务配置。
type AuthServiceConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// NewAuthService 构造认证服务（构造器注入）。
func NewAuthService(users repository.UserRepository, tokenStore TokenStore, cfg AuthServiceConfig) *AuthService {
	return &AuthService{users: users, tokenStore: tokenStore, cfg: cfg, logger: util.Logger}
}

// Login 校验用户名密码并签发 token 对。
func (s *AuthService) Login(ctx context.Context, username, password string) (*model.User, util.TokenPair, error) {
	u, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.TokenPair{}, util.NewAppError(constants.CodeLoginFailed, http.StatusUnauthorized, constants.MessageOf(constants.CodeLoginFailed), err)
		}
		return nil, util.TokenPair{}, fmt.Errorf("auth login: %w", err)
	}
	if !util.CheckPassword(u.PasswordHash, password) {
		return nil, util.TokenPair{}, util.NewAppError(constants.CodeLoginFailed, http.StatusUnauthorized, constants.MessageOf(constants.CodeLoginFailed), nil)
	}
	pair, err := s.issueTokens(ctx, u)
	if err != nil {
		return nil, util.TokenPair{}, err
	}
	return u, pair, nil
}

// Refresh 使用刷新令牌换取新 token 对（轮换）。
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (util.TokenPair, error) {
	userID, err := util.ParseRefreshToken(s.cfg.JWTSecret, refreshToken)
	if err != nil {
		return util.TokenPair{}, util.NewAppError(constants.CodeTokenInvalid, http.StatusUnauthorized, constants.MessageOf(constants.CodeTokenInvalid), err)
	}
	storedID, err := s.tokenStore.GetRefreshTokenUserID(ctx, refreshToken)
	if err != nil {
		return util.TokenPair{}, fmt.Errorf("auth refresh: %w", err)
	}
	if storedID == 0 || storedID != userID {
		return util.TokenPair{}, util.NewAppError(constants.CodeTokenInvalid, http.StatusUnauthorized, constants.MessageOf(constants.CodeTokenInvalid), nil)
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.TokenPair{}, util.NewAppError(constants.CodeUnauthorized, http.StatusUnauthorized, constants.MessageOf(constants.CodeUnauthorized), err)
		}
		return util.TokenPair{}, fmt.Errorf("auth refresh: %w", err)
	}
	if err := s.tokenStore.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return util.TokenPair{}, fmt.Errorf("auth refresh: %w", err)
	}
	return s.issueTokens(ctx, u)
}

// Me 获取当前用户信息。
func (s *AuthService) Me(ctx context.Context, userID uint) (*model.User, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, http.StatusNotFound, constants.MessageOf(constants.CodeUserNotFound), err)
		}
		return nil, fmt.Errorf("auth me: %w", err)
	}
	return u, nil
}

// ChangePassword 修改当前用户密码。
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeUserNotFound, http.StatusNotFound, constants.MessageOf(constants.CodeUserNotFound), err)
		}
		return fmt.Errorf("auth change password: %w", err)
	}
	if !util.CheckPassword(u.PasswordHash, oldPassword) {
		return util.NewAppError(constants.CodePasswordWrong, http.StatusBadRequest, constants.MessageOf(constants.CodePasswordWrong), nil)
	}
	hash, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("auth change password: %w", err)
	}
	u.PasswordHash = hash
	if err := s.users.Update(ctx, u); err != nil {
		return fmt.Errorf("auth change password: %w", err)
	}
	return nil
}

func (s *AuthService) issueTokens(ctx context.Context, u *model.User) (util.TokenPair, error) {
	access, err := util.SignAccessToken(s.cfg.JWTSecret, u.ID, u.Username, string(u.Role), s.cfg.AccessTokenTTL)
	if err != nil {
		return util.TokenPair{}, fmt.Errorf("issue access token: %w", err)
	}
	refresh, err := util.SignRefreshToken(s.cfg.JWTSecret, u.ID)
	if err != nil {
		return util.TokenPair{}, fmt.Errorf("issue refresh token: %w", err)
	}
	if err := s.tokenStore.SetRefreshToken(ctx, refresh, u.ID, s.cfg.RefreshTokenTTL); err != nil {
		return util.TokenPair{}, fmt.Errorf("store refresh token: %w", err)
	}
	return util.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "bearer",
		ExpiresIn:    int64(s.cfg.AccessTokenTTL.Seconds()),
	}, nil
}
