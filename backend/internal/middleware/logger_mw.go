package middleware

import (
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	err := c.Next()
	logger.Log.Info("request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.Int("status", c.Response().StatusCode()),
		zap.String("ip", c.IP()),
	)
	return err
}
