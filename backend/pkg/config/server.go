package config

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServerConfig содержит настройки сервера
type ServerConfig struct {
	Port string
	Env  string
}

// NewServerConfig создает конфигурацию сервера из переменных окружения
func NewServerConfig() *ServerConfig {
	port := GetEnv("SERVER_PORT", "3000")
	env := GetEnv("ENVIRONMENT", "dev")

	return &ServerConfig{
		Port: port,
		Env:  env,
	}
}

// GracefulShutdown выполняет graceful shutdown сервера
func GracefulShutdown(app *fiber.App, db *gorm.DB, timeout time.Duration) {
	logger := zap.L()

	// Создание контекста с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("Начинается graceful shutdown...")

	// Graceful shutdown сервера
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Ошибка при graceful shutdown сервера", zap.Error(err))
	} else {
		logger.Info("Сервер успешно остановлен")
	}

	// Закрытие соединения с базой данных
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("Ошибка при закрытии соединения с БД", zap.Error(err))
			} else {
				logger.Info("Соединение с БД успешно закрыто")
			}
		}
	}

	logger.Info("Graceful shutdown завершен")
}

// ShouldMigrate определяет, нужно ли выполнять миграции
func (c *ServerConfig) ShouldMigrate() bool {
	return c.Env != "dev"
}
