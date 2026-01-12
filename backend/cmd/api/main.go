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

	"github.com/artesipov-alt/odnoi-krovi-app/docsui" // Импорт пакета с обработчиками UI
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/s3"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/seeds"
)

// Options for the CLI.
type Options struct {
	Port int  `help:"Port to listen on" short:"p" default:"3001"`
	Doc  bool `help:"Generate OpenAPI documentation and exit" short:"d"`
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		// Загрузка переменных окружения из .env файла
		godotenv.Load("../.env")

		env := config.GetEnv("ENV", "development")
		miniappDomain := os.Getenv("MINIAPP_DOMAIN") // Получаем домен мини-приложения

		logger.SetupLogger(env)

		// Настройка CORS
		corsHandler := config.SetupCORS(env, miniappDomain)

		// Корневой mux
		rootMux := http.NewServeMux()

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

		if db != nil {
			if err := config.RunMigrations(db); err != nil {
				slog.Error("Ошибка запуска миграций ENT", "error", err)
			}

			ctx := context.Background()
			seeds.SeedBloodGroups(ctx, db)
			seeds.SeedBloodComponents(ctx, db)
			seeds.SeedLocations(ctx, db)
			seeds.SeedBreeds(ctx, db)
		}

		// Инициализация репозиториев
		userRepo := pg.NewEntUserRepository(db)
		locationRepo := pg.NewEntLocationRepository(db)
		breedRepo := pg.NewEntBreedRepository(db)
		bloodInfoRepo := pg.NewEntBloodInfoRepository(db)
		petRepo := pg.NewEntPetRepository(db)
		bloodRequestRepo := pg.NewEntBloodRequestRepository(db)
		fileStorage := s3.NewS3Storage(nil).WithDefaults()

		// Инициализация сервисов
		userService := services.NewUserService(userRepo, locationRepo)
		petService := services.NewPetService(petRepo, userRepo, fileStorage)
		bloodSearchService := services.NewBloodSearchService(bloodRequestRepo, petRepo)

		// Инициализация обработчиков
		userHandler := handlers.NewUserHandler(userService)
		petHandler := handlers.NewPetHandler(petService)
		bloodRequestHandler := handlers.NewBloodRequestHandler(bloodSearchService)
		referenceHandler := handlers.NewReferenceHandler(breedRepo, bloodInfoRepo, locationRepo)

		// Настройка Huma
		api := humago.New(apiMux, config.NewHumaConfig(os.Getenv("MINIAPP_DOMAIN")))

		// Инициализируем интеграцию AppError с Huma
		apperrors.InitHuma(api)

		// Регистрация маршрутов
		userHandler.Register(api)
		petHandler.Register(api)
		bloodRequestHandler.Register(api)
		referenceHandler.Register(api)

		// Если опция Doc включена, генерируем документацию и выходим
		if options.Doc {
			hooks.OnStart(func() {
				generateAndSaveOpenAPI(api, openapiPath)
				os.Exit(0)
			})
			return
		}

		if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				options.Port = port
			} else {
				slog.Warn("Invalid SERVER_PORT environment variable, using default port", "error", err, "value", portStr)
			}
		}

		// Создаем сервер
		server := config.NewServer(options.Port, rootMux)
		// Применяем CORS middleware первым
		server.Use(
			corsHandler.Handler,
			middleware.RecoveryMiddleware,
			middleware.RequestIDMiddleware,
			middleware.LoggingMiddleware,
		)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			slog.Info("Сервер запускается", "port", options.Port)
			server.ListenAndServe()
		})

		// Tell the CLI how to stop your server.
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if db != nil {
				db.Close()
			}
			server.Shutdown(ctx)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}

func generateAndSaveOpenAPI(api huma.API, filePath string) {
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

	slog.Info("OpenAPI documentation successfully generated", "path", filePath)
}
