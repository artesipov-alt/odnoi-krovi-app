// cmd/server/main.go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/handlers"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/middleware"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из .env файла
	godotenv.Load("../.env")

	// Инициализация логгера в режиме разработки
	if os.Getenv("ENV") == "dev" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	}

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

		// ctx := context.Background()
		// // seeds.SeedBloodGroups(ctx, db, slog.Default())
		// // seeds.SeedBloodComponents(ctx, db, slog.Default())
		// // seeds.SeedLocations(ctx, db, slog.Default())
		// // seeds.SeedBreeds(ctx, db, slog.Default())
	}

	// Инициализация репозиториев
	userRepo := pg.NewEntUserRepository(db)
	locationRepo := pg.NewEntLocationRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo, locationRepo)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)

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

	// Создаем сервер и вешаем middleware на сырой HTTP (mux)
	server := config.NewServer(mux)

	// Теперь Use принимает func(http.Handler) http.Handler
	server.Use(
		middleware.RecoveryMiddleware,
		middleware.RequestHandler,
	)

	// Канал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Сервер запускается", "port", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Ошибка запуска сервера", "error", err)
			os.Exit(1)
		}
	}()

	<-quit

	slog.Info("🚨 Получен сигнал завершения работы сервера")
	server.Shutdown(context.Background())
}
