package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		method := c.Method()
		path := c.Path()

		// 记录请求开始
		zap.L().Debug("-> 请求开始",
			zap.String("method", method),
			zap.String("path", path),
		)

		// 处理请求
		err := c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 记录请求完成
		zap.L().Info("<- 请求完成",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.IP()),
		)
		return err
	}
}
