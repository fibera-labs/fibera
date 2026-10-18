package common

import (
	"github.com/gofiber/fiber/v3"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 返回成功响应
func Success(c fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Code:    fiber.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
