// cmd/server/main.go
package main

import (
	"context"
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
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/s3"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/validator"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/seeds"

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
		logger.SetupSlogDefaultLogger(env)

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
		db, err := config.ConnectEnt(config.NewENVConfig())
		if err != nil {
			slog.Error("Ошибка подключения к базе данных (ENT)", "error", err)
			os.Exit(1)
		}

		// Запуск миграций
		if err := config.RunMigrations(db); err != nil {
			slog.Error("Ошибка выполнения миграций", "error", err)
			os.Exit(1)
		}

		//Миграции
		ctx := context.Background()
		seeds.SeedBloodGroups(ctx, db)
		seeds.SeedBloodComponents(ctx, db)
		seeds.SeedLocations(ctx, db)
		seeds.SeedBreeds(ctx, db)

		// Инициализация репозиториев
		userRepo := pg.NewEntUserRepository(db)
		locationRepo := pg.NewEntLocationRepository(db)
		breedRepo := pg.NewEntBreedRepository(db)
		bloodInfoRepo := pg.NewEntBloodInfoRepository(db)
		petRepo := pg.NewEntPetRepository(db)
		bloodRequestRepo := pg.NewEntBloodRequestRepository(db)
		fileStorage := s3.NewS3Storage(nil).WithDefaults()

		donorValidator := validator.NewDonorValidator(validator.DefaultStopChecks, validator.DefaultWarnChecks)

		// Инициализация сервисов
		userService := services.NewUserService(userRepo, locationRepo, fileStorage)
		petService := services.NewPetService(petRepo, userRepo, bloodRequestRepo, fileStorage, donorValidator)
		bloodSearchService := services.NewBloodSearchService(bloodRequestRepo, petRepo, fileStorage, db)
		fileService := services.NewFileService(petRepo, userRepo, bloodRequestRepo, fileStorage)

		// Инициализация обработчиков
		referenceHandler := handlers.NewReferenceHandler(breedRepo, bloodInfoRepo, locationRepo)
		userHandler := handlers.NewUserHandler(userService, locationRepo)
		petHandler := handlers.NewPetHandler(petService, donorValidator)
		bloodRequestHandler := handlers.NewBloodRequestHandler(bloodSearchService)
		fileHandler := handlers.NewFileHandler(fileService)

		// Настройка Huma
		humapi = humago.New(apiMux, config.NewHumaConfig(os.Getenv("MINIAPP_DOMAIN")))

		// Инициализируем интеграцию AppError с Huma
		apperrors.InitHuma(humapi)

		// Регистрация маршрутов
		userHandler.Register(humapi)
		petHandler.Register(humapi)
		bloodRequestHandler.Register(humapi)
		fileHandler.Register(humapi)
		referenceHandler.Register(humapi)

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
		server.Use(
			sloghttp.Recovery,
			sloghttp.New(slog.Default()),
			config.DefaultCorsHandler(env, miniappDomain),
		)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			slog.Info("Сервер запускается", "port", options.Port)
			server.ListenAndServe()
		})

		// Tell the CLI how to stop your server.
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if db != nil {
				db.Close()
			}
			server.Shutdown(ctx)
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
