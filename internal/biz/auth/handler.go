package auth

import (
	"errors"
	"strings"

	"fibera/internal/common"

	"github.com/gofiber/fiber/v3"
)

// Handler 认证请求处理器
type Handler struct {
	service *Service
}

// NewHandler 创建认证请求处理器
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login 登录认证
//
//	@Summary		用户登录
//	@Description	用户ID与用户名二选一（同时提供时以用户ID为准），登录成功返回 token 与用户信息
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"登录请求参数"
//	@Success		200		{object}	common.Response{data=auth.LoginResponse}
//	@Failure		400		{object}	common.Response	"用户名（或用户ID）和密码不能为空"
//	@Failure		401		{object}	common.Response	"用户名或密码错误"
//	@Failure		403		{object}	common.Response	"用户已被禁用"
//	@Failure		500		{object}	common.Response	"服务内部错误"
//	@Router			/api/auth/login [post]
func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.Error(c, fiber.StatusBadRequest, "请求参数有误")
	}
	hasUserID := req.UserID != nil && *req.UserID > 0
	if req.Password == "" || (!hasUserID && strings.TrimSpace(req.UserName) == "") {
		return common.Error(c, fiber.StatusBadRequest, "用户名（或用户ID）和密码不能为空")
	}

	result, err := h.service.Login(&req, c.IP())
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return common.Error(c, fiber.StatusUnauthorized, err.Error())
		case errors.Is(err, ErrUserDisabled):
			return common.Error(c, fiber.StatusForbidden, err.Error())
		default:
			return common.Error(c, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.Success(c, result)
}
