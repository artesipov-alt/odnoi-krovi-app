package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
)

// EntConfig содержит конфигурацию для подключения к PostgreSQL через ENT
type EntConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewEntConfig создает конфигурацию из переменных окружения
func NewEntConfig() *EntConfig {
	return &EntConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "odnoi_krovi"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// GetDSN возвращает строку подключения для PostgreSQL
func (c *EntConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// ConnectEnt подключается к PostgreSQL и возвращает экземпляр ent.Client
func ConnectEnt() (*ent.Client, error) {
	config := NewEntConfig()
	dsn := config.GetDSN()

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	return client, nil
}

// RunMigrations запускает автоматическую миграцию схем
func RunMigrations(client *ent.Client) error {
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	log.Println("ENT migrations completed successfully")
	return nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
