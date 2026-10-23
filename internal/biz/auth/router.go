package auth

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes 注册认证模块路由
func RegisterRoutes(app *fiber.App) {
	handler := NewHandler(NewService())
	group := app.Group("/api/auth")

	group.Post("/login", handler.Login)
}
