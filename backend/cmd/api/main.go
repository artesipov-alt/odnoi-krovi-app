// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/joho/godotenv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/s3"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/config"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/seeds"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"3001"`
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Загрузка переменных окружения из .env файла
		godotenv.Load("../.env")

		logger.SetupLogger(config.GetEnv("ENV", "development"))

		// Инициализация кастомных ошибок для Huma
		// Это переопределяет huma.NewError, чтобы использовать ваш AppError
		apperrors.InitHuma()

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
		devHandler := handlers.NewDevHandler(userRepo)
		bloodRequestHandler := handlers.NewBloodRequestHandler(bloodSearchService)

		// Создание стандартного mux
		mux := http.NewServeMux()

		// Настройка Huma
		humaConfig := huma.DefaultConfig("1krovi.app API", "1.3.5")
		// Мы больше не переопределяем humaConfig.ErrorHandler, так как
		// переопределили huma.NewError в apperrors.InitHuma().
		// Huma сама будет создавать AppError и записывать его в Response.

		api := humago.New(mux, humaConfig)

		// Регистрация маршрутов
		userHandler.Register(api)
		petHandler.Register(api)
		devHandler.Register(api)
		bloodRequestHandler.Register(api)

		// Создаем сервер
		server := &http.Server{Addr: fmt.Sprintf(":%d", options.Port), Handler: mux}

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			slog.Info("Сервер запускается", "port", options.Port)
			server.ListenAndServe()
		})

		// Tell the CLI how to stop your server.
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
