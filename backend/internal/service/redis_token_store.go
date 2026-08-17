// Package service 实现全部业务逻辑。
package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenStore 刷新令牌存储接口。
type TokenStore interface {
	SetRefreshToken(ctx context.Context, token string, userID uint, ttl time.Duration) error
	GetRefreshTokenUserID(ctx context.Context, token string) (uint, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

// RedisTokenStore 基于 Redis 的刷新令牌存储实现。
type RedisTokenStore struct {
	client *redis.Client
}

// NewRedisTokenStore 构造 Redis 令牌存储。
func NewRedisTokenStore(client *redis.Client) *RedisTokenStore {
	return &RedisTokenStore{client: client}
}

const refreshTokenKeyPrefix = "supplychain:refresh:"

func (s *RedisTokenStore) SetRefreshToken(ctx context.Context, token string, userID uint, ttl time.Duration) error {
	if err := s.client.Set(ctx, refreshTokenKeyPrefix+token, userID, ttl).Err(); err != nil {
		return fmt.Errorf("set refresh token: %w", err)
	}
	return nil
}

func (s *RedisTokenStore) GetRefreshTokenUserID(ctx context.Context, token string) (uint, error) {
	v, err := s.client.Get(ctx, refreshTokenKeyPrefix+token).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("get refresh token: %w", err)
	}
	uid, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse refresh token user id: %w", err)
	}
	return uint(uid), nil
}

func (s *RedisTokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	if err := s.client.Del(ctx, refreshTokenKeyPrefix+token).Err(); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}
