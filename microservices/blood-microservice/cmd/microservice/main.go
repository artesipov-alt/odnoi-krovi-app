package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	bloodsearchv1connect "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1/bloodsearchv1connect"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/repositories/pg"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/pkg/logger"
	"github.com/joho/godotenv"

	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	// Инициализируем логгер
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	if err := logger.Init(env); err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Создаем Chi роутер
	r := chi.NewRouter()
	r.Use(middleware.RequestID,
		middleware.ZapLogger,
		chimiddleware.Recoverer,
		chimiddleware.RealIP)

	// Подключаемся к PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=bloodsearch port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Не удалось подключиться к PostgreSQL", zap.Error(err))
	}

	// Автомиграция таблиц
	if err := db.AutoMigrate(&models.PetRow{}); err != nil {
		logger.Log.Fatal("Не удалось выполнить миграцию таблиц", zap.Error(err))
	}

	logger.Log.Info("Успешное подключение к PostgreSQL")

	// Создаем репозиторий
	bloodSearchRepo := pg.NewBloodSearchRepositoryGorm(db)

	// Создаем сервис
	bloodSearchService := services.NewBloodSearchService(bloodSearchRepo, logger.Log)

	// Создаем Connect handler
	path, handler := bloodsearchv1connect.NewBloodSearchPoolHandler(bloodSearchService)

	// Подключаем handler к Chi
	r.Handle(path+"*", handler)

	// Настраиваем протоколы для поддержки HTTP/2 без TLS
	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	// Создаем HTTP сервер с поддержкой HTTP/2
	s := http.Server{
		Addr:      ":8081",
		Handler:   r,
		Protocols: p,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		logger.Log.Info("Сервер запущен на :8081")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	// Ожидаем сигналы для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Получен сигнал завершения работы...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Ошибка при завершении работы сервера", zap.Error(err))
	}

	logger.Log.Info("Сервер корректно завершил работу")
}
