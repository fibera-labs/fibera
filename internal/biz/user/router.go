package user

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes 注册用户模块路由
func RegisterRoutes(app *fiber.App) {
	handler := NewHandler(NewService())
	group := app.Group("/api/users")

	group.Get("/", handler.GetUserList)
	group.Get("/:id", handler.GetUserByID)
	group.Post("/", handler.CreateUser)
	group.Put("/:id", handler.UpdateUser)
	group.Delete("/:id", handler.DeleteUser)
}
