package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 生成密码哈希。
func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", fmt.Errorf("hash password: empty password")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(h), nil
}

// CheckPassword 校验明文密码与哈希是否匹配。
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
