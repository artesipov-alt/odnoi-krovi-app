// cmd/server/main.go
package main

import (
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func main() {
	logger.Init("dev")
	defer logger.Sync()

	app := fiber.New()

	app.Use(func(c fiber.Ctx) error {
		err := c.Next()
		logger.Log.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("ip", c.IP()),
		)
		return err
	})

	app.Get("/", func(c fiber.Ctx) error {
		logger.Log.Info("root accessed")
		return c.SendString("Hello from Fiber + Zap 👋")
	})

	app.Listen(":3000")
}
