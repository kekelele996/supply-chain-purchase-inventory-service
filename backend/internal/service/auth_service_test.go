package service

import (
	"context"
	"testing"
	"time"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

func newAuthServiceForTest(t *testing.T) (*AuthService, *fakeUserRepo) {
	t.Helper()
	users := newFakeUserRepo()
	_ = users.Create(context.Background(), mustUser(t, "admin", "admin123", string(constants.RoleAdmin)))
	svc := NewAuthService(users, newFakeTokenStore(), AuthServiceConfig{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  24 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
	return svc, users
}

func TestAuthService_LoginSuccess(t *testing.T) {
	svc, _ := newAuthServiceForTest(t)
	user, pair, err := svc.Login(context.Background(), "admin", "admin123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if user.Username != "admin" {
		t.Fatalf("unexpected user: %+v", user)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatalf("tokens must not be empty")
	}
	claims, err := util.ParseAccessToken("test-secret", pair.AccessToken)
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	if claims.Username != "admin" || claims.Role != string(constants.RoleAdmin) {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestAuthService_LoginWrongPassword(t *testing.T) {
	svc, _ := newAuthServiceForTest(t)
	_, _, err := svc.Login(context.Background(), "admin", "wrongpass")
	if err == nil {
		t.Fatalf("expected error")
	}
	ae := util.AsAppError(err)
	if ae == nil || ae.Code != constants.CodeLoginFailed {
		t.Fatalf("expected login failed error, got %v", err)
	}
}

func TestAuthService_RefreshRotation(t *testing.T) {
	svc, _ := newAuthServiceForTest(t)
	_, pair, err := svc.Login(context.Background(), "admin", "admin123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	newPair, err := svc.Refresh(context.Background(), pair.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if newPair.AccessToken == "" {
		t.Fatalf("new access token must not be empty")
	}
	// 旧 refresh token 已轮换，再次使用应失败
	if _, err := svc.Refresh(context.Background(), pair.RefreshToken); err == nil {
		t.Fatalf("expected refresh with rotated token to fail")
	}
}

func TestAuthService_RefreshInvalidToken(t *testing.T) {
	svc, _ := newAuthServiceForTest(t)
	if _, err := svc.Refresh(context.Background(), "not-a-real-token"); err == nil {
		t.Fatalf("expected error for invalid refresh token")
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	svc, users := newAuthServiceForTest(t)
	u, _ := users.FindByUsername(context.Background(), "admin")
	if err := svc.ChangePassword(context.Background(), u.ID, "admin123", "newpass123"); err != nil {
		t.Fatalf("change password: %v", err)
	}
	if err := svc.ChangePassword(context.Background(), u.ID, "wrong", "newpass123"); err == nil {
		t.Fatalf("expected error for wrong old password")
	}
	// 新密码可登录
	if _, _, err := svc.Login(context.Background(), "admin", "newpass123"); err != nil {
		t.Fatalf("login with new password: %v", err)
	}
}
