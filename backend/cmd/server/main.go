// cmd/server/main.go
package main

import (
	_ "github.com/artesipov-alt/odnoi-krovi-app/docs"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title однойкрови.рф
// @version 1.0
// @description API сервиса однойкрови.рф для донороcства крови и помощи животным
// @Host 212.111.85.99:3000
// @BasePath /api/v1
func main() {
	// Загрузка переменных окружения из .env файла
	godotenv.Load()

	logger.Init("dev")
	defer logger.Sync()

	app := fiber.New()

	app.Use(cors.New(config.CORSOptions()))

	app.Use(middleware.LoggerMiddleware)

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes group
	api := app.Group("/api/v1")
	{
		api.Get("/", handlers.RootHandler)

		api.Get("/user/:id", handlers.GetUserHandler)
		api.Post("/user/", handlers.AddUserHandler)
		// Здесь будут добавляться другие API endpoints
	}

	app.Listen(":3000")
}
