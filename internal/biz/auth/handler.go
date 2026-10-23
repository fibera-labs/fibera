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
