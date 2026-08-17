package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/dto"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

// UserHandler 用户管理处理器（admin 专用）。
type UserHandler struct {
	users  *service.UserService
	logger *slog.Logger
}

// NewUserHandler 构造用户管理处理器。
func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users, logger: util.Logger}
}

// List 用户列表。
func (h *UserHandler) List(c *gin.Context) {
	page := util.ParsePageQuery(c)
	username := c.Query("username")
	role := c.Query("role")
	list, total, err := h.users.List(c.Request.Context(), page.Page, page.PageSize, username, role)
	if err != nil {
		util.Error(c, err)
		return
	}
	items := make([]dto.UserResponse, 0, len(list))
	for _, u := range list {
		items = append(items, toUserResponse(u))
	}
	util.OK(c, dto.PageResult{List: items, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Create 创建用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	u, err := h.users.Create(c.Request.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Created(c, constants.MsgUserCreated, toUserResponse(*u))
}

// Get 用户详情。
func (h *UserHandler) Get(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	u, err := h.users.Get(c.Request.Context(), idReq.ID)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OK(c, toUserResponse(*u))
}

// Update 更新用户。
func (h *UserHandler) Update(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	u, err := h.users.Update(c.Request.Context(), idReq.ID, req.Password, req.Role)
	if err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUserUpdated, toUserResponse(*u))
}

// Delete 删除用户。
func (h *UserHandler) Delete(c *gin.Context) {
	var idReq dto.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		util.Error(c, util.BuildValidationError(err))
		return
	}
	if err := h.users.Delete(c.Request.Context(), idReq.ID); err != nil {
		util.Error(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUserDeleted, nil)
}
