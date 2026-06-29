// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/artesipov-alt/odnoi-krovi-app/docsui" // Импорт пакета с обработчиками UI
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/analytics/query"
	authcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/auth/cmd"
	bloodcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/cmd"
	bloodquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/query"
	bonuscmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bonus/cmd"
	donorcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/donor/cmd"
	donorquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/donor/query"
	filecmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/file/cmd"
	petcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/cmd"
	petquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/query"
	refquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/reference/query"
	usercmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/cmd"
	userquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/user/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/otp/twin24"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/pg"
	redisRepository "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/redis"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/s3"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/scheduler"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/scheduler/job"

	events "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/events/redis"
	transport "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/middleware"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/auth"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"

	sloghttp "github.com/samber/slog-http"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"3001"`
}

func main() {
	var humapi huma.API

	// Создаем CLI инструмент сервера.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Загрузка переменных окружения из .env файла
		godotenv.Load("../.env")
		env := config.GetEnv("ENV", "development")
		miniappDomain := os.Getenv("MINIAPP_DOMAIN") // Получаем домен мини-приложения

		// Заменяем стандартный слог логером от Charm Bracelet.
		logger.SetupSlogDefaultLogger(os.Getenv("LOG_LEVEL"))

		// Корневой mux
		rootMux := http.NewServeMux()

		// Health check endpoint
		rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// API mux с префиксом /api
		apiMux := http.NewServeMux()

		// Подключаем API mux к /api
		rootMux.Handle("/api/", http.StripPrefix("/api", apiMux))

		//Указываем директорию документации
		openapiPath := filepath.Join("docs", "openapi.json")

		// Обслуживание файла openapi.json
		apiMux.Handle("/openapi.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, openapiPath)
		}))

		// Обслуживание UI документации Swagger
		apiMux.HandleFunc("/docs", docsui.ScalarDocsHandler)

		// Инициализация подключения к базе данных через ENT
		db, rawdb, err := config.ConnectEnt(config.NewEntConfig(env))
		if err != nil {
			slog.Error("Ошибка подключения к базе данных (ENT)", "error", err)
			os.Exit(1)
		}

		var otpRepo redisRepository.OTPRepository
		var publisher ports.EventPublisher
		var otpSender *twin24.OTPSender
		redisClient, err := config.NewRedisClientFromEnv()
		if err != nil {
			slog.Warn("Redis недоступен, события не будут публиковаться", "error", err)
			publisher = &events.NoOpEventPublisher{}
			otpRepo = redisRepository.NewNoOpOTPRepo()
		} else {
			publisher = events.NewEventPublisher(redisClient, env)
			otpRepo = redisRepository.NewOTPRepo(redisClient)
			otpSender = twin24.NewOTPSenderFromEnv(redisClient)
		}

		// Запуск миграций закомментирован, так как они больше не нужны.
		if err := config.RunMigrations(db, rawdb); err != nil {
			slog.Error("Ошибка выполнения миграций", "error", err)
			os.Exit(1)
		}

		//Миграции
		// ctx := context.Background()
		// seeds.SeedLocations(ctx, db)
		// seeds.SeedBreeds(ctx, db)

		// Инициализация репозиториев
		userRepo := pg.NewEntUserRepository(db)
		locationRepo := pg.NewEntLocationRepository(db)
		breedRepo := pg.NewEntBreedRepository(db)
		petRepo := pg.NewEntPetRepository(db)
		bloodRequestRepo := pg.NewEntBloodRequestRepository(db)
		donorResponseRepo := pg.NewEntDonorResponseRepository(db)
		partnerRepo := pg.NewEntPartnerRepository(db)
		bonusRepo := pg.NewEntBonusRepository(db)
		rawQueryRepo := pg.NewRawQueryRepository(rawdb)
		bonusSvc := bonus.NewBonusService(bonusRepo)
		fileStorage := s3.NewS3Storage(nil).WithDefaults()
		txManager := presistance.NewTxManager(db)

		// Инициализация reference query handlers
		getAllBreedsHandler := refquery.NewGetAllBreedsHandler(breedRepo)
		getBreedsByTypeHandler := refquery.NewGetBreedsByPetTypeHandler(breedRepo)
		getAllLocationsHandler := refquery.NewGetAllLocationsHandler(locationRepo)
		getAllBloodComponentsHandler := refquery.NewGetAllBloodComponentsHandler()
		getBloodGroupsByTypeHandler := refquery.NewGetBloodGroupsByPetTypeHandler()

		getPortalStatisticsHandler := query.NewPortalStatsHandler(rawQueryRepo)

		//Дополнительные сервисы для аунтификации
		// miniAppDataValidator := auth.NewAppValidator(os.Getenv("TG_BOT_TOKEN"), os.Getenv("MAX_BOT_TOKEN"))
		miniAppDataValidator := auth.NewSimpleValidator()
		tokenGenerator := auth.NewJWTGenerator(os.Getenv("JWT_SECRET_KEY"), "odnoi-krovi-backend", 24*time.Hour)

		externalSignInHandler := authcmd.NewExternalSignInHandler(userRepo, partnerRepo, tokenGenerator, txManager)
		appSgnInHandler := authcmd.NewMiniAppSignInHandler(userRepo, miniAppDataValidator, tokenGenerator, txManager)

		userDeleteHandler := usercmd.NewDeleteHandler(userRepo)
		userUpdateHandler := usercmd.NewUpdateHandler(userRepo, txManager)
		userChangePhoneHandler := usercmd.NewChangePhoneHandler(userRepo, otpRepo, otpSender)
		userVerifyPhoneHandler := usercmd.NewVerifyPhoneHandler(userRepo, otpRepo, txManager)
		userResetHandler := usercmd.NewResetHandler(userRepo)
		userRestoreHandler := usercmd.NewRestoreHandler(userRepo)
		userGetByIDHandler := userquery.NewGetByIDHandler(userRepo)
		userGetContactHandler := userquery.NewGetContactHandler(userRepo, publisher)
		userGetDeletedHandler := userquery.NewGetDeletedUsersHandler(userRepo)

		matchingSvc := *bloodsearch.NewMatchingService()
		petService := pet.NewPetService()

		donorGetRecipientsListHandler := donorquery.NewListRequestsHandler(petRepo, donorResponseRepo, bloodRequestRepo, matchingSvc, petService, userRepo)
		donorApplyBloodHandler := donorcmd.NewApplyForRequestHandler(bloodRequestRepo, petRepo, donorResponseRepo, userRepo, bonusSvc, publisher, txManager)
		donorGetRecipientDetailsHandler := donorquery.NewRecipientDetailHandler(donorResponseRepo, petRepo, bloodRequestRepo, userRepo, matchingSvc, petService, bonusSvc)
		donorGetPlannedDonationsHandler := donorquery.NewPlannedDonationsHandler(donorResponseRepo, petRepo, bloodRequestRepo, userRepo, bonusRepo)
		donorGetCompletedDonationsHandler := donorquery.NewCompletedDonationsHandler(donorResponseRepo, petRepo, bloodRequestRepo, userRepo, bonusRepo)
		assignedBonusesHandler := donorquery.NewAssignedBonusesHandler(userRepo, bonusRepo)
		completeDonationHandler := donorcmd.NewCompleteDonationHandler(donorResponseRepo, bloodRequestRepo, petRepo, userRepo, publisher)
		cancelDonationHandler := donorcmd.NewCancelDonationHandler(donorResponseRepo, bloodRequestRepo, txManager, publisher, petRepo, userRepo, bonusSvc)

		// Инициализация pet handlers
		// petRepo реализует все интерфейсы: PetReadRepository, PetWriteRepository, PetStatsRepository, PetPhotoRepository
		petCreateHandler := petcmd.NewCreateHandler(petRepo, userRepo)
		petUpdateHandler := petcmd.NewUpdateHandler(petRepo, petRepo)
		petDeleteHandler := petcmd.NewDeleteHandler(petRepo, petRepo, bloodRequestRepo)
		petGetByIDHandler := petquery.NewGetByIDHandler(petRepo, bloodRequestRepo)
		petGetByUserHandler := petquery.NewGetByUserHandler(petRepo, userRepo, donorResponseRepo, bloodRequestRepo, bonusRepo, petService)

		// Инициализация bloodsearch handlers
		bloodCreateHandler := bloodcmd.NewCreateRequestHandler(bloodRequestRepo, petRepo, donorResponseRepo, userRepo, bonusRepo, publisher, petService, txManager)
		bloodUpdateHandler := bloodcmd.NewUpdateRequestHandler(bloodRequestRepo)
		bloodDeleteHandler := bloodcmd.NewDeleteRequestHandler(bloodRequestRepo, txManager)
		bloodGetByIDHandler := bloodquery.NewGetByIDHandler(bloodRequestRepo, petRepo)
		bloodGetByPetIDHandler := bloodquery.NewGetByPetIDHandler(bloodRequestRepo, petRepo)
		bloodGetDonorByIDHandler := bloodquery.NewGetDonorByIDHandler(petRepo, donorResponseRepo, bloodRequestRepo, petService)
		bloodGetDonationHandler := bloodquery.NewGetDonationHandler(petRepo, donorResponseRepo, userRepo, bloodRequestRepo, petService)
		applyResponseHandler := bloodcmd.NewApplyResponseHandler(bloodRequestRepo, donorResponseRepo, petRepo, userRepo, publisher, txManager)
		confirmDonationHandler := bloodcmd.NewConfirmDonationHandler(bloodRequestRepo, donorResponseRepo, petRepo, userRepo, txManager, publisher, bonusSvc)
		bloodCloseDonationHandler := bloodcmd.NewCloseRequestHandler(bloodRequestRepo, donorResponseRepo, petRepo, userRepo, txManager, publisher, bonusSvc)
		rejectDonationHandler := bloodcmd.NewRejectDonationHandler(bloodRequestRepo, donorResponseRepo, petRepo, userRepo, txManager, publisher, bonusSvc)
		// Инициализация file handlers
		fileGetPresignedHandler := filecmd.NewGetPresignedURLsHandler(fileStorage, petRepo, userRepo, bloodRequestRepo)
		fileConfirmUploadHandler := filecmd.NewConfirmUploadHandler(petRepo, userRepo, bloodRequestRepo, fileStorage)
		bonusImportHandler := bonuscmd.NewImportBonusesHandler(bonusRepo)

		// Инициализация handlers
		authHandler := transport.NewAuthHandler(
			appSgnInHandler,
			externalSignInHandler,
		)

		referenceHandler := transport.NewReferenceHandler(
			getAllBreedsHandler,
			getBreedsByTypeHandler,
			getAllLocationsHandler,
			getAllBloodComponentsHandler,
			getBloodGroupsByTypeHandler,
		)
		userHandler := transport.NewUserHandler(
			userDeleteHandler,
			userUpdateHandler,
			userChangePhoneHandler,
			userVerifyPhoneHandler,
			userResetHandler,
			userRestoreHandler,
			userGetByIDHandler,
			userGetContactHandler,
			userGetDeletedHandler,
			fileStorage,
		)
		petHandler := transport.NewPetHandler(
			petCreateHandler,
			petUpdateHandler,
			petDeleteHandler,
			petGetByIDHandler,
			petGetByUserHandler,
			fileStorage,
		)
		bloodRequestHandler := transport.NewBloodRequestHandler(
			bloodCreateHandler,
			bloodUpdateHandler,
			bloodDeleteHandler,
			bloodGetByIDHandler,
			bloodGetByPetIDHandler,
			bloodGetDonorByIDHandler,
			bloodGetDonationHandler,
			applyResponseHandler,
			confirmDonationHandler,
			rejectDonationHandler,
			bloodCloseDonationHandler,
			fileStorage,
		)
		donorHandler := transport.NewDonorHandler(
			donorGetRecipientsListHandler,
			donorGetRecipientDetailsHandler,
			donorApplyBloodHandler,
			donorGetPlannedDonationsHandler,
			donorGetCompletedDonationsHandler,
			completeDonationHandler,
			cancelDonationHandler,
			assignedBonusesHandler,
			fileStorage,
		)

		fileHandler := transport.NewFileHandler(
			fileGetPresignedHandler,
			fileConfirmUploadHandler,
			bonusImportHandler,
		)

		commonHandler := transport.NewCommonHandler(
			getPortalStatisticsHandler,
		)

		//Запуск side-effects (воркеров)
		confirmJob := job.NewAutoConfirmJob(donorResponseRepo, confirmDonationHandler)
		scheduler := scheduler.NewScheduler()
		scheduler.Register(confirmJob, 10*time.Minute)
		scheduler.Start()

		// Настройка Huma
		humapi = humago.New(apiMux, config.NewHumaConfig(os.Getenv("MINIAPP_DOMAIN")))

		// Инициализируем интеграцию AppError с Huma
		apperrors.InitHuma(humapi)

		// Регистрация маршрутов
		authHandler.Register(humapi)
		userHandler.Register(humapi)
		petHandler.Register(humapi)
		donorHandler.Register(humapi)
		bloodRequestHandler.Register(humapi)
		fileHandler.Register(humapi)
		referenceHandler.Register(humapi)
		commonHandler.Register(humapi)

		if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				options.Port = port
			} else {
				slog.Warn("Invalid SERVER_PORT environment variable, using default port", "error", err, "value", portStr)
			}
		}

		// Создаем сервер с корневым mux
		server := config.NewServer(options.Port, rootMux)

		// Применяем middleware с использованием метода Use
		// Порядок: Recovery -> CORS -> BasicAuth -> Auth -> Logging -> Mux
		server.Use(
			sloghttp.Recovery,
			config.DefaultCorsHandler(env, miniappDomain),
			middleware.BasicAuthMiddleware("/api/docs", "/api/openapi.json"),
			middleware.AuthMiddleware(tokenGenerator, env, "/api/v1/auth", "/api/docs", "/api/openapi.json", "/health"),
			sloghttp.NewWithConfig(slog.Default(), sloghttp.Config{WithResponseBody: true, WithRequestID: true, Filters: []sloghttp.Filter{sloghttp.AcceptStatusGreaterThanOrEqual(400)}}),
			middleware.TraceIDResponseMiddleware,
		)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			slog.Info("✅ Сервер Запустился", "url", fmt.Sprintf("http://localhost:%d/api/docs", options.Port))
			server.ListenAndServe()
		})

		// Tell the CLI how to stop your server.
		hooks.OnStop(func() {
			slog.Info("🛑 Завершение работы сервера...")

			// 1. Останавливаем scheduler (воркеры)
			scheduler.Stop()

			// 2. Graceful shutdown HTTP-сервера (30s timeout)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				slog.Error("Ошибка при завершении HTTP-сервера", "error", err)
			}

			// 3. Закрываем подключение к БД (Ent)
			if db != nil {
				if err := db.Close(); err != nil {
					slog.Error("Ошибка при закрытии БД", "error", err)
				}
			}

			// 4. Закрываем Redis-клиент
			if redisClient != nil {
				if err := redisClient.Close(); err != nil {
					slog.Error("Ошибка при закрытии Redis", "error", err)
				}
			}

			slog.Info("✅ Сервер успешно остановлен")
		})
	})

	// Добавляем команду для генерации документации.
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Generate the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			GenerateOpenAPI(humapi, "./docs/openapi.json")
			slog.Info("Спецификация создана!")
		},
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}

func GenerateOpenAPI(api huma.API, filePath string) {
	// Создаем директорию, если она не существует
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("Failed to create directory for OpenAPI spec", "error", err, "path", dir)
		return
	}

	spec := api.OpenAPI()
	jsonBytes, err := spec.MarshalJSON()
	if err != nil {
		slog.Error("Failed to marshal OpenAPI spec to JSON", "error", err)
		return
	}

	err = os.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		slog.Error("Failed to write OpenAPI spec to file", "error", err, "path", filePath)
		return
	}
}
