// cmd/server/main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/artesipov-alt/odnoi-krovi-app/docs" // Документация Swagger
	cache "github.com/artesipov-alt/odnoi-krovi-app/internal/cache/redis"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers" // Обработчики HTTP запросов
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware"

	// Промежуточное ПО
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/s3"

	// Репозитории для работы с БД
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	cacherepo "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/cache"
	pgrepositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg" // Репозитории для работы с БД
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"                       // Бизнес-логика
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/config"                   // Конфигурация приложения
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"                   // Логирование
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/seeds"                    // Сиды для БД

	// Управление миграциями
	"github.com/joho/godotenv"                              // Загрузка .env файлов
	"github.com/labstack/echo/v4"                           // Веб-фреймворк
	echomiddleware "github.com/labstack/echo/v4/middleware" // CORS middleware
	echoSwagger "github.com/swaggo/echo-swagger"            // Swagger UI
	"go.uber.org/zap"                                       // Структурированное логирование
)

// @title 1krovi.app
// @version 1.3.5
// @description API сервиса однойкрови.рф для донороcства крови и помощи животным
// @openapi 3.0.0
// @servers https://1krovi.app {description: "Production server", url: https://1krovi.app/api}
// @servers http://localhost:3001 {description: "Local development server", url: http://localhost:3001/api}
// @BasePath /api
func main() {

	// Загрузка переменных окружения из .env файла
	godotenv.Load("../.env")

	// Инициализация конфигурации сервера
	serverConfig := config.NewServerConfig()

	// Инициализация логгера в режиме разработки
	logger.Init(serverConfig.Env)
	defer logger.Sync() // Гарантированное закрытие логгера при завершении

	// Инициализация подключения к базе данных через ENT в зависимости от окружения
	db, err := config.ConnectEnt(config.NewENVConfig())
	if err != nil {
		logger.Log.Fatal("Ошибка подключения к базе данных (ENT)", zap.Error(err))
	}

	// Запуск миграций ENT (если необходимо)
	if serverConfig.ShouldMigrate(true) {
		if err := config.RunMigrations(db); err != nil {
			logger.Log.Error("Ошибка запуска миграций ENT", zap.Error(err))
		}

		// Заполнение БД начальными данными (Seeds)
		ctx := context.Background()
		seeds.SeedBloodGroups(ctx, db, logger.Log)
		seeds.SeedBloodComponents(ctx, db, logger.Log)
		seeds.SeedLocations(ctx, db, logger.Log)
		seeds.SeedBreeds(ctx, db, logger.Log)
	}

	// Создаем кэш, но если ошибка - используем nil
	rCache, err := cache.NewCacheFromEnv()
	if err != nil {
		logger.Log.Warn("Redis недоступен, работаем без кэша", zap.Error(err))
		rCache = nil
	}

	// Инициализация репозиториев
	userRepo := pgrepositories.NewEntUserRepository(db)
	petRepo := pgrepositories.NewEntPetRepository(db)
	breedRepo := pgrepositories.NewEntBreedRepository(db)
	bloodRepo := pgrepositories.NewEntBloodInfoRepository(db)
	locationRepo := pgrepositories.NewEntLocationRepository(db)
	bloodreqRepo := pgrepositories.NewEntBloodRequestRepository(db)

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
	userService := services.NewUserService(userRepo, locationRepo)
	petService := services.NewPetService(petRepo, userRepo, s3VKCloud)
	bloodReqService := services.NewBloodSearchService(bloodreqRepo, petRepo)

	// Инициализация обработчиков HTTP запросов (хэндлеров)
	userHandler := handlers.NewUserHandler(userService)
	petHandler := handlers.NewPetHandler(petService)
	bloodRequestHandler := handlers.NewBloodRequestHandler(bloodReqService)
	referenceHandler := handlers.NewReferenceHandler(breedRepo, bloodRepoInit, locationRepo)
	devHandler := handlers.NewDevHandler(userRepo)

	// Создание экземпляра echo приложения с кастомным обработчиком ошибок
	app := echo.New()
	app.HTTPErrorHandler = middleware.ErrorHandler()

	// Настройка CORS для кросс-доменных запросов
	app.Use(echomiddleware.CORSWithConfig(echomiddleware.DefaultCORSConfig))

	// Подключение middleware
	app.Use(echomiddleware.Recover())   // Восстановление после паники
	app.Use(middleware.RequestLogger()) // Логирование запросов с кастомным логгером

	// Группировка API маршрутов с префиксом /api
	api := app.Group("/api")

	// Документация Swagger - доступна по адресу /api/swagger/index.html
	api.GET("/swagger/*", echoSwagger.WrapHandlerV3)

	// Группа маршрутов для версии API v1
	v1 := api.Group("/v1")

	// Группа маршрутов для работы с пользователями
	userGroup := v1.Group("/user")

	userGroup.GET("/telegram", userHandler.GetUserByTelegramHandler)          // Получение пользователя по Telegram ID
	userGroup.POST("/register", userHandler.RegisterUserHandler)              // Регистрация нового пользователя
	userGroup.POST("/register/simple", userHandler.RegisterUserSimpleHandler) // Простая регистрация (для команды Start)
	userGroup.GET("/:id", userHandler.GetUserHandler)                         // Получение пользователя по ID
	userGroup.PUT("/:id", userHandler.UpdateUserHandler)                      // Обновление данных пользователя
	userGroup.DELETE("/:id", userHandler.DeleteUserHandler)                   // Удаление пользователя по ID

	// Группа маршрутов для разработчиков
	devGroup := v1.Group("/dev")

	devGroup.POST("/restore-user/:id", devHandler.RestoreUserHandler)
	devGroup.POST("/reset-user/:id", devHandler.ResetUserHandler)     // Сброс пользователя к заводским настройкам
	devGroup.GET("/deleted-users", devHandler.GetDeletedUsersHandler) // Получение всех удаленных пользователей

	// Группа маршрутов для работы с питомцами и поиском крови
	petGroup := v1.Group("/pet")

	petGroup.GET("/user/:user_id", petHandler.GetUserPetsHandler)                    // Получение всех питомцев пользователя
	petGroup.POST("/user/:user_id", petHandler.CreatePetHandler)                     // Создание питомца для пользователя
	petGroup.GET("/:id", petHandler.GetPetHandler)                                   // Получение питомца по ID
	petGroup.PUT("/:id", petHandler.UpdatePetHandler)                                // Обновление данных питомца
	petGroup.DELETE("/:id", petHandler.DeletePetHandler)                             // Удаление питомца по ID
	petGroup.GET("/upload/avatar/:id", petHandler.GetAvatarUploadURL)                // Получение ссылки на загрузку в фотографии питомцев в storage
	petGroup.POST("/upload/avatar/confirm/:path", petHandler.ConfirmPetAvatarUpload) // Подтверждение загрузки

	// Группа маршрутов для пула запросов крови
	bloodRequestGroup := v1.Group("/blood-request")

	bloodRequestGroup.POST("/pool", bloodRequestHandler.AddPetToBloodRequestPool)           // Добавить питомца в пул поиска крови
	bloodRequestGroup.POST("/pool/search", bloodRequestHandler.GetPetsFromBloodRequestPool) // Получить питомцев из пула поиска крови

	// Группа маршрутов для справочных данных
	referenceGroup := v1.Group("/reference")

	referenceGroup.GET("/pet-types", referenceHandler.GetPetTypesHandler)
	referenceGroup.GET("/pet-roles", referenceHandler.GetPetRolesHandler)
	referenceGroup.GET("/genders", referenceHandler.GetGendersHandler)
	referenceGroup.GET("/user-roles", referenceHandler.GetUserRolesHandler)
	referenceGroup.GET("/breeds", referenceHandler.GetBreedsHandler)
	referenceGroup.GET("/breeds-by-type", referenceHandler.GetBreedsByTypeHandler)
	referenceGroup.GET("/blood-components", referenceHandler.GetBloodComponentsHandler)
	referenceGroup.GET("/blood-groups/:pet_type", referenceHandler.GetBloodGroupsHandler)
	referenceGroup.GET("/living-conditions", referenceHandler.GetLivingConditionsHandler)
	referenceGroup.GET("/locations", referenceHandler.GetLocationsHandler)
	referenceGroup.GET("/health-statuses", referenceHandler.GetHealthStatusesHandler)
	referenceGroup.GET("/reproductive-statuses", referenceHandler.GetReproductiveStatusesHandler)

	// Канал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Log.Info("Сервер запускается", zap.String("port", serverConfig.Port))
		if err := app.Start(":" + serverConfig.Port); err != nil {
			logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-quit

	// Подробное логирование перед graceful shutdown
	logger.Log.Info("🚨 Получен сигнал завершения работы сервера")
	config.GracefulShutdown(app, db, rCache, 30)
}
