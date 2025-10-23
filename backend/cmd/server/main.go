// cmd/server/main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/artesipov-alt/odnoi-krovi-app/docs"              // Документация Swagger
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"   // Обработчики HTTP запросов
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware" // Промежуточное ПО
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"

	// Модели данных
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
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	env := os.Getenv("ENVIROMENT")
	if env == "" {
		env = "dev"
	}
	// Инициализация логгера в режиме разработки
	logger.Init(env)
	defer logger.Sync() // Гарантированное закрытие логгера при завершении

	// Инициализация подключения к базе данных
	db, err := config.ConnectDB()
	if err != nil {
		logger.Log.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}

	// Автоматическое создание/обновление таблиц в БД
	// autoMigrate(db)

	// Инициализация репозитория для работы с пользователями
	userRepo := repositories.NewPostgresUserRepository(db)

	// Инициализация репозитория для работы с питомцами
	petRepo := repositories.NewPostgresPetRepository(db)

	// Инициализация репозитория для работы с породами
	breedRepo := repositories.NewPostgresBreedRepository(db)

	// Инициализация репозитория для работы с типами крови
	bloodTypeRepo := repositories.NewPostgresBloodTypeRepository(db)

	// Инициализация сервиса с бизнес-логикой пользователей
	userService := services.NewUserService(userRepo)

	// Инициализация сервиса с бизнес-логикой питомцев
	petService := services.NewPetService(petRepo, userRepo)

	// Инициализация обработчиков HTTP запросов для пользователей
	userHandler := handlers.NewUserHandler(userService)

	// Инициализация обработчиков HTTP запросов для питомцев
	petHandler := handlers.NewPetHandler(petService)

	// Инициализация обработчиков для справочных данных
	referenceHandler := handlers.NewReferenceHandler(breedRepo, bloodTypeRepo)

	// Создание экземпляра Fiber приложения
	app := fiber.New()

	// Настройка CORS для кросс-доменных запросов
	app.Use(cors.New(config.CORSOptions()))

	// Подключение middleware
	app.Use(middleware.LoggerMiddleware) // Логирование запросов
	// app.Use(middleware.TelegramAuthMiddleware(middleware.DefaultTelegramAuthConfig())) // Реальная аутентификация Telegram (закомментирована)
	// app.Use(middleware.MockTelegramAuthMiddleware(middleware.DefaultMockTelegramConfig())) // Тестовая аутентификация Telegram

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
			userGroup.Delete("/:id", userHandler.DeleteUserHandler)                   // Удаление пользователя по ID
		}

		// Группа маршрутов для работы с питомцами
		petGroup := api.Group("/pets")
		{
			petGroup.Get("/user/:user_id", petHandler.GetUserPetsHandler) // Получение всех питомцев пользователя
			petGroup.Post("/user/:user_id", petHandler.CreatePetHandler)  // Создание питомца для пользователя
			petGroup.Get("/:id", petHandler.GetPetHandler)                // Получение питомца по ID
			petGroup.Put("/:id", petHandler.UpdatePetHandler)             // Обновление данных питомца
			petGroup.Delete("/:id", petHandler.DeletePetHandler)          // Удаление питомца по ID
		}

		// Группа маршрутов для справочных данных
		referenceGroup := api.Group("/reference")
		{
			referenceGroup.Get("/pet-types", referenceHandler.GetPetTypesHandler)                 // Типы животных
			referenceGroup.Get("/genders", referenceHandler.GetGendersHandler)                    // Пол животного
			referenceGroup.Get("/living-conditions", referenceHandler.GetLivingConditionsHandler) // Условия проживания
			referenceGroup.Get("/user-roles", referenceHandler.GetUserRolesHandler)               // Роли пользователей
			referenceGroup.Get("/breeds", referenceHandler.GetBreedsHandler)
			referenceGroup.Get("/breeds-by-type", referenceHandler.GetBreedsByTypeHandler)               // Породы животных
			referenceGroup.Get("/blood-groups", referenceHandler.GetBloodGroupsHandler)                  // Группы крови
			referenceGroup.Get("/blood-search-statuses", referenceHandler.GetBloodSearchStatusesHandler) // Статусы поиска крови
			referenceGroup.Get("/blood-stock-statuses", referenceHandler.GetBloodStockStatusesHandler)   // Статусы запаса крови
			referenceGroup.Get("/donation-statuses", referenceHandler.GetDonationStatusesHandler)        // Статусы донорства
		}
	}

	// Канал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Log.Info("Сервер запускается", zap.String("port", port))
		if err := app.Listen(":" + port); err != nil {
			logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	logger.Log.Info("Получен сигнал завершения, начинается graceful shutdown...")

	// Создание контекста с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown сервера
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Log.Error("Ошибка при graceful shutdown сервера", zap.Error(err))
	}

	// Закрытие соединения с базой данных
	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Log.Error("Ошибка при закрытии соединения с БД", zap.Error(err))
		} else {
			logger.Log.Info("Соединение с БД успешно закрыто")
		}
	}

	logger.Log.Info("Сервер успешно остановлен")
}

// Автоматическое создание/обновление таблиц в базе данных
func autoMigrate(db *gorm.DB) {
	modelsToMigrate := []any{
		&models.User{},
		&models.Pet{},
		&models.VetClinic{},
		&models.Breed{},
		&models.BloodSearch{},
		&models.BloodStock{},
		&models.BloodType{},
		&models.Location{},
		&models.Donation{},
	}

	// Автоматическая миграция всех моделей
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		logger.Log.Fatal("Ошибка автоматической миграции", zap.Error(err))
	}
}
