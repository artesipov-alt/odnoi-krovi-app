package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// Root godoc
// @Summary Корневой эндпоинт
// @Description Возвращает информационное сообщение о сервере
// @Tags root
// @Produce plain
// @Success 200 {string} string "Информационное сообщение"
// @Router / [get]
func RootHandler(c *fiber.Ctx) error {
	logger.Log.Info("root accessed")
	return c.SendString("Это тестовый бэкенд сервер на Go + Fiber + Swagger проекта однойкрови.рф\nДокументация API доступна по адресу: /swagger/")
}
