// cmd/server/main.go
package main

import (
	_ "github.com/artesipov-alt/odnoi-krovi-app/docs"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// @title однойкрови.рф
// @version 1.0
// @description API сервиса однойкрови.рф для донороcства крови и помощи животным
// @host
// @BasePath /api/v1
func main() {
	// Загрузка переменных окружения из .env файла
	godotenv.Load()

	logger.Init("dev")
	defer logger.Sync()

	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate database models
	autoMigrate(db)

	// Initialize repositories
	userRepo := repositories.NewPostgresUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	app := fiber.New()

	app.Use(cors.New(config.CORSOptions()))

	app.Use(middleware.LoggerMiddleware)
	// app.Use(middleware.TelegramAuthMiddleware(middleware.DefaultTelegramAuthConfig()))
	app.Use(middleware.MockTelegramAuthMiddleware(middleware.DefaultMockTelegramConfig()))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes group
	api := app.Group("/api/v1")
	{
		api.Get("/", handlers.RootHandler)

		// User routes
		userGroup := api.Group("/user")
		{
			userGroup.Get("/telegram", userHandler.GetUserByTelegramHandler)
			userGroup.Post("/register", userHandler.RegisterUserHandler)
			userGroup.Get("/:id", userHandler.GetUserHandler)
			userGroup.Put("/:id", userHandler.UpdateUserHandler)
		}
	}

	logger.Log.Info("Server starting on :3000")
	if err := app.Listen(":3000"); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}

// Auto-migrate database models
func autoMigrate(db *gorm.DB) {
	// TODO: Add models for auto-migration
	db.AutoMigrate(&models.User{})
}
