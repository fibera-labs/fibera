package app

import (
	"fibera/internal/biz/user"
	"fibera/internal/config"
	"fibera/internal/infrastructure/db"
	"fibera/internal/infrastructure/log"
	customMiddleware "fibera/internal/middleware"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// Start 启动应用程序
func Start() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load configuration file: %v", err)
		return
	}

	err = setupInfrastructure(conf)
	if err != nil {
		fmt.Printf("failed to setup infrastructure: %v", err)
		return
	}

	app := fiber.New()
	registerMiddleware(app)
	registerRoutes(app)

	err = app.Listen(":" + conf.App.Port)
	if err != nil {
		return
	}
}

// setupInfrastructure 初始化应用基础设施服务
func setupInfrastructure(conf *config.Config) (err error) {
	// 初始化日志
	err = log.SetupLogger(conf)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %v", err)
	}
	// 初始化数据库连接
	if _, err = db.GetDB(); err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	// 同步数据库表结构
	if conf.Database.Sync {
		err = db.SyncDatabase()
		if err != nil {
			return fmt.Errorf("failed to sync database: %v", err)
		}
	}
	return nil
}

// registerMiddleware 注册中间件
func registerMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(customMiddleware.Logger())
	app.Use(customMiddleware.Jwt())
}

// registerRoutes 注册路由
func registerRoutes(app *fiber.App) {
	user.RegisterRoutes(app)
}
