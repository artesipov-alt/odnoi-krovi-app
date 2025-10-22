// cmd/server/main.go
package main

import (
	_ "github.com/artesipov-alt/odnoi-krovi-app/docs"                // Документация Swagger
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"     // Обработчики HTTP запросов
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware"   // Промежуточное ПО
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"       // Модели данных
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories" // Репозитории для работы с БД
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"     // Бизнес-логика
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"            // Конфигурация приложения
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"            // Логирование
	"github.com/gofiber/fiber/v2"                                    // Веб-фреймворк
	"github.com/gofiber/fiber/v2/middleware/cors"                    // CORS middleware
	"github.com/gofiber/swagger"                                     // Swagger UI
	"github.com/joho/godotenv"                                       // Загрузка .env файлов
	"go.uber.org/zap"                                                // Структурированное логирование
	"gorm.io/gorm"                                                   // ORM для работы с БД
)

// @title однойкрови.рф
// @version 1.0
// @description API сервиса однойкрови.рф для донороcства крови и помощи животным
// @host
// @BasePath /api/v1
func main() {
	// Загрузка переменных окружения из .env файла
	godotenv.Load()

	// Инициализация логгера в режиме разработки
	logger.Init("dev")
	defer logger.Sync() // Гарантированное закрытие логгера при завершении

	// Инициализация подключения к базе данных
	db, err := config.ConnectDB()
	if err != nil {
		logger.Log.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}

	// Автоматическое создание/обновление таблиц в БД
	autoMigrate(db)

	// Инициализация репозитория для работы с пользователями
	userRepo := repositories.NewPostgresUserRepository(db)

	// Инициализация сервиса с бизнес-логикой пользователей
	userService := services.NewUserService(userRepo)

	// Инициализация обработчиков HTTP запросов для пользователей
	userHandler := handlers.NewUserHandler(userService)

	// Создание экземпляра Fiber приложения
	app := fiber.New()

	// Настройка CORS для кросс-доменных запросов
	app.Use(cors.New(config.CORSOptions()))

	// Подключение middleware
	app.Use(middleware.LoggerMiddleware) // Логирование запросов
	// app.Use(middleware.TelegramAuthMiddleware(middleware.DefaultTelegramAuthConfig())) // Реальная аутентификация Telegram (закомментирована)
	app.Use(middleware.MockTelegramAuthMiddleware(middleware.DefaultMockTelegramConfig())) // Тестовая аутентификация Telegram

	// Документация Swagger - доступна по адресу /swagger/*
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Группировка API маршрутов с префиксом /api/v1
	api := app.Group("/api/v1")
	{
		// Корневой маршрут API
		api.Get("/", handlers.RootHandler)

		// Группа маршрутов для работы с пользователями
		userGroup := api.Group("/user")
		{
			userGroup.Get("/telegram", userHandler.GetUserByTelegramHandler)          // Получение пользователя по Telegram ID
			userGroup.Post("/register", userHandler.RegisterUserHandler)              // Регистрация нового пользователя
			userGroup.Post("/register/simple", userHandler.RegisterUserSimpleHandler) // Простая регистрация (для команды Start)
			userGroup.Get("/:id", userHandler.GetUserHandler)                         // Получение пользователя по ID
			userGroup.Put("/:id", userHandler.UpdateUserHandler)                      // Обновление данных пользователя
		}
	}

	// Запуск сервера на порту 3000
	logger.Log.Info("Сервер запускается на :3000")
	if err := app.Listen(":3000"); err != nil {
		logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}

// Автоматическое создание/обновление таблиц в базе данных
func autoMigrate(db *gorm.DB) {
	// TODO: Добавить другие модели для автоматической миграции
	db.AutoMigrate(&models.User{}) // Создание таблицы пользователей
}
