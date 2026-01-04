package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/labstack/echo/v4"
)

// Root godoc
// @Summary Корневой эндпоинт
// @Description Возвращает информационное сообщение о сервере
// @Tags root
// @Produce html
// @Success 200 {string} string "Информационное сообщение с HTML-ссылкой"
// @Router / [get]
func RootHandler(c echo.Context) error {
	logger.Log.Info("root accessed")
	htmlResponse := `<html>
<head>
				<title>Одной Крови</title>
				<style>
								body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
								h1 { color: #333; }
								p { font-size: 18px; }
								a { color: #007BFF; text-decoration: none; }
								a:hover { text-decoration: underline; }
				</style>
</head>
<body>
				<h1>Добро пожаловать!</h1>
				<p>Это тестовый бэкенд сервер на Go + Echo + Swagger проекта однойкрови.рф</p>
				<p>Документация API доступна по адресу: <a href="https://1krovi.app/api/swagger/">https://1krovi.app/api/swagger/</a></p>
</body>
</html>`
	return c.HTML(200, htmlResponse)
}
