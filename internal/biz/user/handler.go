package user

import (
	"errors"
	"strconv"

	"fibera/internal/common"
	"fibera/internal/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Handler 系统用户请求处理器
type Handler struct {
	service *Service
}

// NewHandler 创建系统用户请求处理器
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetUserList 获取用户列表
func (h *Handler) GetUserList(c fiber.Ctx) error {
	users, err := h.service.GetUserList()
	if err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, users)
}

// GetUserByID 获取用户详情
func (h *Handler) GetUserByID(c fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return common.Error(c, fiber.StatusBadRequest, "无效的用户ID")
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Error(c, fiber.StatusNotFound, "用户不存在")
		}
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, user)
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c fiber.Ctx) error {
	var user model.SysUser
	if err := c.Bind().Body(&user); err != nil {
		return common.Error(c, fiber.StatusBadRequest, "请求参数有误")
	}

	if err := h.service.CreateUser(&user); err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, user)
}

// UpdateUser 更新用户
func (h *Handler) UpdateUser(c fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return common.Error(c, fiber.StatusBadRequest, "无效的用户ID")
	}

	var user model.SysUser
	if err := c.Bind().Body(&user); err != nil {
		return common.Error(c, fiber.StatusBadRequest, "请求参数有误")
	}
	user.ID = id

	if err := h.service.UpdateUser(&user); err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, user)
}

// DeleteUser 删除用户
func (h *Handler) DeleteUser(c fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return common.Error(c, fiber.StatusBadRequest, "无效的用户ID")
	}

	if err := h.service.DeleteUser(id); err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, nil)
}

// parseUserID 解析路径中的用户主键
func parseUserID(c fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
