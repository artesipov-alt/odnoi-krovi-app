package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodpool/v1/bloodpoolv1connect"
	redisrepo "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/repositories/redis"
	v1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	redisclient "github.com/redis/go-redis/v9"
)

func main() {
	// Создаем Chi роутер
	r := chi.NewRouter()
	r.Use(middleware.Logger,
		middleware.Recoverer,
		middleware.RealIP)

	// Создаем Redis клиент
	redisClient := redisclient.NewClient(&redisclient.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Проверяем подключение к Redis
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	log.Println("Успешное подключение к Redis")

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
					log.Printf("failed to clean expired pool: %v", err)
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
		log.Println("Сервер запущен на :8080")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигналы для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Получен сигнал завершения работы...")

	// Останавливаем background cleaner
	cleanerCancel()

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}

	log.Println("Сервер корректно завершил работу")
}
