// cmd/server/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/artesipov-alt/odnoi-krovi-app/docs" // Документация Swagger
	cache "github.com/artesipov-alt/odnoi-krovi-app/internal/cache/redis"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"   // Обработчики HTTP запросов
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware" // Промежуточное ПО
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/s3"

	// Репозитории для работы с БД
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	cacherepo "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/cache"
	pgrepositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg" // Репозитории для работы с БД
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"                       // Бизнес-логика
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"                              // Конфигурация приложения
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"                              // Логирование
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/migration"

	// Управление миграциями
	"github.com/gofiber/fiber/v2"                 // Веб-фреймворк
	"github.com/gofiber/fiber/v2/middleware/cors" // CORS middleware
	"github.com/gofiber/swagger"                  // Swagger UI
	"github.com/joho/godotenv"                    // Загрузка .env файлов
	"go.uber.org/zap"                             // Структурированное логирование
)

// @title 1krovi.app
// @version 1.3.1
// @description API сервиса однойкрови.рф для донороcства крови и помощи животным
// @host
// @BasePath /api/v1
func main() {

	// Загрузка переменных окружения из .env файла
	godotenv.Load("../.env")

	// Инициализация конфигурации сервера
	serverConfig := config.NewServerConfig()

	// Инициализация логгера в режиме разработки
	logger.Init(serverConfig.Env)
	defer logger.Sync() // Гарантированное закрытие логгера при завершении

	// Инициализация подключения к базе данных
	db, err := config.ConnectDB()
	if err != nil {
		logger.Log.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}

	// Создаем кэш, но если ошибка - используем nil
	rCache, err := cache.NewCacheFromEnv()
	if err != nil {
		logger.Log.Warn("Redis недоступен, работаем без кэша", zap.Error(err))
		rCache = nil
	}

	// Автоматическое создание/обновление таблиц в БД на проде
	migration.AutoMigrate(db, logger.Log)
	migration.SeedDatabase(db, logger.Log)

	// Инициализация репозиториев
	userRepo := pgrepositories.NewPostgresUserRepository(db)
	petRepo := pgrepositories.NewPostgresPetRepository(db)
	breedRepo := pgrepositories.NewPostgresBreedRepository(db)
	bloodRepo := pgrepositories.NewPostgresBloodRepository(db)
	locationRepo := pgrepositories.NewPostgresLocationRepository(db)

	// Создаем репозиторий в зависимости от наличия кэша
	var bloodRepoInit repositories.BloodInfoRepository
	if rCache != nil {
		// Кеширующие репозитории
		bloodRepoInit = cacherepo.NewCachedBloodInfoRepository(bloodRepo, rCache)
	} else {
		// Используем обычный репозиторий без кэша
		bloodRepoInit = bloodRepo
	}
	// Инициализация S3
	s3VKCloud := s3.NewS3Storage(nil).WithDefaults()

	// Инициализация сервисов
	userService := services.NewUserService(userRepo)
	petService := services.NewPetService(petRepo, userRepo, s3VKCloud)

	// Инициализация клиента blood search микросервиса
	bloodSearchClient := *services.NewBloodRequestClient(serverConfig.BloodMicroserviceURL)

	// Инициализация обработчиков HTTP запросов (хэндлеров)
	userHandler := handlers.NewUserHandler(userService)
	petHandler := handlers.NewPetHandler(petService, bloodSearchClient)
	referenceHandler := handlers.NewReferenceHandler(breedRepo, bloodRepoInit, locationRepo)
	devHandler := handlers.NewDevHandler(userRepo)

	// Создание экземпляра Fiber приложения с кастомным обработчиком ошибок
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Настройка CORS для кросс-доменных запросов
	app.Use(cors.New(config.CORSOptions()))

	// Подключение middleware
	app.Use(middleware.RecoveryMiddleware()) // Восстановление после паники
	app.Use(middleware.LoggerMiddleware)     // Логирование запросов
	// app.Use(middleware.TelegramAuthMiddleware(middleware.DefaultTelegramAuthConfig())) // Реальная аутентификация Telegram (закомментирована)
	// app.Use(middleware.MockTelegramAuthMiddleware(middleware.DefaultMockTelegramConfig())) // Тестовая аутентификация Telegram

	// Группировка API маршрутов с префиксом /api
	api := app.Group("/api")

	// Документация Swagger - доступна по адресу /api/swagger/*
	api.Get("/swagger/*", swagger.HandlerDefault)

	// Группировка API маршрутов с префиксом /api/v1
	v1 := api.Group("/v1")

	// Корневой маршрут API
	v1.Get("/", handlers.RootHandler)

	// Группа маршрутов для работы с пользователями
	userGroup := v1.Group("/user")

	userGroup.Get("/telegram", userHandler.GetUserByTelegramHandler)          // Получение пользователя по Telegram ID
	userGroup.Post("/register", userHandler.RegisterUserHandler)              // Регистрация нового пользователя
	userGroup.Post("/register/simple", userHandler.RegisterUserSimpleHandler) // Простая регистрация (для команды Start)
	userGroup.Get("/:id", userHandler.GetUserHandler)                         // Получение пользователя по ID
	userGroup.Put("/:id", userHandler.UpdateUserHandler)                      // Обновление данных пользователя
	userGroup.Delete("/:id", userHandler.DeleteUserHandler)                   // Удаление пользователя по ID

	// Группа маршрутов для разработчиков
	devGroup := v1.Group("/dev")

	devGroup.Post("/restore-user/:id", devHandler.RestoreUserHandler)
	devGroup.Post("/reset-user/:id", devHandler.ResetUserHandler)     // Сброс пользователя к заводским настройкам
	devGroup.Get("/deleted-users", devHandler.GetDeletedUsersHandler) // Получение всех удаленных пользователей

	// Группа маршрутов для работы с питомцами и поиском крови
	petGroup := v1.Group("/pets")

	petGroup.Get("/user/:user_id", petHandler.GetUserPetsHandler)                   // Получение всех питомцев пользователя
	petGroup.Post("/user/:user_id", petHandler.CreatePetHandler)                    // Создание питомца для пользователя
	petGroup.Get("/:id", petHandler.GetPetHandler)                                  // Получение питомца по ID
	petGroup.Put("/:id", petHandler.UpdatePetHandler)                               // Обновление данных питомца
	petGroup.Delete("/:id", petHandler.DeletePetHandler)                            // Удаление питомца по ID
	petGroup.Get("upload/avatar/:id", petHandler.GetAvatarUploadURL)                // Получение ссылки на загрузку в фотографии питомцев в storage
	petGroup.Post("upload/avatar/confirm/:path", petHandler.ConfirmPetAvatarUpload) // Подтверждение загрузки

	// Поиск крови связан с питомцами: добавление и поиск питомцев для поиска крови
	petGroup.Post("/blood-request/pool", petHandler.AddPetToBloodRequestPool)           // Добавить питомца в пул поиска крови
	petGroup.Post("/blood-request/pool/search", petHandler.GetPetsFromBloodRequestPool) // Получить питомцев из пула поиска крови

	// Группа маршрутов для справочных данных
	referenceGroup := v1.Group("/reference")

	referenceGroup.Get("/pet-types", referenceHandler.GetPetTypesHandler)
	referenceGroup.Get("/pet-roles", referenceHandler.GetPetRolesHandler)
	referenceGroup.Get("/genders", referenceHandler.GetGendersHandler)
	referenceGroup.Get("/user-roles", referenceHandler.GetUserRolesHandler)
	referenceGroup.Get("/breeds", referenceHandler.GetBreedsHandler)
	referenceGroup.Get("/breeds-by-type", referenceHandler.GetBreedsByTypeHandler)
	referenceGroup.Get("/blood-components", referenceHandler.GetBloodComponentsHandler)
	referenceGroup.Get("/blood-groups/:pet_type", referenceHandler.GetBloodGroupsHandler)
	referenceGroup.Get("/living-conditions", referenceHandler.GetLivingConditionsHandler)
	referenceGroup.Get("/locations", referenceHandler.GetLocationsHandler)
	referenceGroup.Get("/health-statuses", referenceHandler.GetHealthStatusesHandler)
	referenceGroup.Get("/reproductive-statuses", referenceHandler.GetReproductiveStatusesHandler)

	// Канал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Log.Info("Сервер запускается", zap.String("port", serverConfig.Port))
		if err := app.Listen(":" + serverConfig.Port); err != nil {
			logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-quit

	// Подробное логирование перед graceful shutdown
	logger.Log.Info("🚨 Получен сигнал завершения работы сервера")

	// Graceful shutdown сервера
	config.GracefulShutdown(app, db, rCache, 30*time.Second)
}
