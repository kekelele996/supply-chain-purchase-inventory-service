package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair 访问令牌与刷新令牌。
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

var (
	// ErrInvalidToken token 无效。
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken token 过期。
	ErrExpiredToken = errors.New("expired token")
)

// SignAccessToken 签发访问令牌。
func SignAccessToken(secret string, userID uint, username, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "supplychain-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return signed, nil
}

// SignRefreshToken 签发刷新令牌（附带 token_type=refresh 标记与唯一 jti）。
func SignRefreshToken(secret string, userID uint) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "refresh",
		"jti":        NewRequestID(),
		"iss":        "supplychain-api",
		"iat":        now.Unix(),
		"exp":        now.Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}
	return signed, nil
}

// ParseAccessToken 校验并解析访问令牌。
func ParseAccessToken(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == 0 {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// ParseRefreshToken 校验刷新令牌并返回 user_id。
func ParseRefreshToken(secret, tokenString string) (uint, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrExpiredToken
		}
		return 0, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !token.Valid {
		return 0, ErrInvalidToken
	}
	if tt, ok := claims["token_type"].(string); !ok || tt != "refresh" {
		return 0, ErrInvalidToken
	}
	uid, ok := claims["user_id"].(float64)
	if !ok || uid <= 0 {
		return 0, ErrInvalidToken
	}
	return uint(uid), nil
}
