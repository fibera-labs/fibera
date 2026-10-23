package user

import (
	"errors"
	"strconv"

	"fibera/internal/common"

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
//
//	@Summary		获取用户列表
//	@Description	查询全部系统用户，返回结果不包含密码等敏感字段
//	@Tags			系统用户
//	@Produce		json
//	@Success		200	{object}	common.Response{data=[]user.UserResponse}
//	@Failure		401	{object}	common.Response	"未认证或 token 无效"
//	@Failure		500	{object}	common.Response	"服务内部错误"
//	@Security		BearerAuth
//	@Router			/api/users [get]
func (h *Handler) GetUserList(c fiber.Ctx) error {
	users, err := h.service.GetUserList()
	if err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, users)
}

// GetUserByID 获取用户详情
//
//	@Summary		获取用户详情
//	@Description	按主键查询单个系统用户，返回结果不包含密码等敏感字段
//	@Tags			系统用户
//	@Produce		json
//	@Param			id	path		int	true	"用户主键"
//	@Success		200	{object}	common.Response{data=user.UserResponse}
//	@Failure		400	{object}	common.Response	"无效的用户ID"
//	@Failure		401	{object}	common.Response	"未认证或 token 无效"
//	@Failure		404	{object}	common.Response	"用户不存在"
//	@Failure		500	{object}	common.Response	"服务内部错误"
//	@Security		BearerAuth
//	@Router			/api/users/{id} [get]
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
//
//	@Summary		创建用户
//	@Description	创建系统用户，密码在落库前进行加密处理
//	@Tags			系统用户
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateUserRequest	true	"创建用户请求参数"
//	@Success		200		{object}	common.Response{data=user.UserResponse}
//	@Failure		400		{object}	common.Response	"请求参数有误"
//	@Failure		401		{object}	common.Response	"未认证或 token 无效"
//	@Failure		500		{object}	common.Response	"服务内部错误"
//	@Security		BearerAuth
//	@Router			/api/users [post]
func (h *Handler) CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.Error(c, fiber.StatusBadRequest, "请求参数有误")
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, user)
}

// UpdateUser 更新用户
//
//	@Summary		更新用户
//	@Description	按主键更新系统用户信息
//	@Tags			系统用户
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"用户主键"
//	@Param			request	body		UpdateUserRequest	true	"更新用户请求参数"
//	@Success		200		{object}	common.Response{data=user.UserResponse}
//	@Failure		400		{object}	common.Response	"无效的用户ID或请求参数有误"
//	@Failure		401		{object}	common.Response	"未认证或 token 无效"
//	@Failure		500		{object}	common.Response	"服务内部错误"
//	@Security		BearerAuth
//	@Router			/api/users/{id} [put]
func (h *Handler) UpdateUser(c fiber.Ctx) error {
	id, err := parseUserID(c)
	if err != nil {
		return common.Error(c, fiber.StatusBadRequest, "无效的用户ID")
	}

	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.Error(c, fiber.StatusBadRequest, "请求参数有误")
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		return common.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return common.Success(c, user)
}

// DeleteUser 删除用户
//
//	@Summary		删除用户
//	@Description	按主键删除系统用户
//	@Tags			系统用户
//	@Produce		json
//	@Param			id	path		int	true	"用户主键"
//	@Success		200	{object}	common.Response	"删除成功"
//	@Failure		400	{object}	common.Response	"无效的用户ID"
//	@Failure		401	{object}	common.Response	"未认证或 token 无效"
//	@Failure		500	{object}	common.Response	"服务内部错误"
//	@Security		BearerAuth
//	@Router			/api/users/{id} [delete]
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
