package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// Root godoc
// @Summary Корневой эндпоинт
// @Description Возвращает информационное сообщение о сервере
// @Tags root
// @Produce html
// @Success 200 {string} string "Информационное сообщение с HTML-ссылкой"
// @Router / [get]
func RootHandler(c *fiber.Ctx) error {
	logger.Log.Info("root accessed")
	htmlResponse := `Это тестовый бэкенд сервер на Go + Fiber + Swagger проекта однойкрови.рф<br>
Документация API доступна по адресу: <a href="https://1krovi.app/api/swagger/">https://1krovi.app/api/swagger/</a>`
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(htmlResponse)
}
