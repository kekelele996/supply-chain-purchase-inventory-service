// Package handler 实现 HTTP 层请求处理。
package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// AuthHandler 认证处理器。
type AuthHandler struct {
	auth   *service.AuthService
	logger *slog.Logger
}

// NewAuthHandler 构造认证处理器。
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth, logger: util.Logger}
}

// Login 用户登录。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user, pair, err := h.auth.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Warn("auth login failed", "username", req.Username, "err", err)
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgLoginSuccess, dto.TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    pair.TokenType,
		ExpiresIn:    pair.ExpiresIn,
	})
	_ = user
}

// Refresh 刷新 token。
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	pair, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgRefreshSuccess, dto.TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    pair.TokenType,
		ExpiresIn:    pair.ExpiresIn,
	})
}

// Me 获取当前用户信息。
func (h *AuthHandler) Me(c *gin.Context) {
	user := util.CurrentUser(c)
	u, err := h.auth.Me(c.Request.Context(), user.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, dto.UserProfileResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
	})
}

// ChangePassword 修改密码。
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	user := util.CurrentUser(c)
	if err := h.auth.ChangePassword(c.Request.Context(), user.ID, req.OldPassword, req.NewPassword); err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgPasswordChanged, nil)
}

// Health 健康检查。
func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": constants.MsgHealthOK})
}
