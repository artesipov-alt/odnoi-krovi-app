package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodpool/v1/bloodpoolv1connect"
	redisrepo "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/repositories/redis"
	v1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/pkg/logger"

	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
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

	// Создаем Redis клиент
	redisClient := redisclient.NewClient(&redisclient.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Проверяем подключение к Redis
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("Не удалось подключиться к Redis", zap.Error(err))
	}
	logger.Log.Info("Успешное подключение к Redis")

	// Создаем репозиторий
	petRepo := redisrepo.NewRedisPetRepository(redisClient)

	// Background cleaner: периодически очищает истекшие записи из ZSET-пула.
	cleanerCtx, cleanerCancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := petRepo.CleanExpiredPool(cleanerCtx); err != nil {
					logger.Log.Error("failed to clean expired pool", zap.Error(err))
				}
			case <-cleanerCtx.Done():
				return
			}
		}
	}()

	// Создаем gRPC сервер
	poolService := v1.NewBloodPoolService(petRepo)
	path, handler := bloodpoolv1connect.NewBloodSearchPoolHandler(poolService)

	// Подключаем gRPC handler к Chi
	r.Handle(path+"*", handler)

	// Настраиваем протоколы для поддержки HTTP/2 без TLS
	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	// Создаем HTTP сервер с поддержкой HTTP/2
	s := http.Server{
		Addr:      ":8080",
		Handler:   r,
		Protocols: p,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		logger.Log.Info("Сервер запущен на :8080")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	// Ожидаем сигналы для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Получен сигнал завершения работы...")

	// Останавливаем background cleaner
	cleanerCancel()

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Ошибка при завершении работы сервера", zap.Error(err))
	}

	logger.Log.Info("Сервер корректно завершил работу")
}
