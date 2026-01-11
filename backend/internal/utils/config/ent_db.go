package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	_ "github.com/artesipov-alt/odnoi-krovi-app/ent/runtime"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/pg"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
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
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", "postgres"),
		DBName:   GetEnv("DB_NAME", "odnoi_krovi"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}
}

// NewLocalConfig создает конфигурацию для локальной разработки
func NewLocalConfig() *EntConfig {
	return &EntConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "adminpass1921",
		DBName:   "local_odnoi_krovi",
		SSLMode:  "disable",
	}
}

func NewENVConfig() *EntConfig {
	if os.Getenv("ENV") != "development" {
		return NewEntConfig()
	}
	return NewLocalConfig()
}

// GetDSN возвращает строку подключения для PostgreSQL
func (c *EntConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// ConnectEnt подключается к PostgreSQL и возвращает экземпляр ent.Client
func ConnectEnt(config *EntConfig) (*ent.Client, error) {
	dsn := config.GetDSN()

	// Открываем соединение через стандартный sql.DB
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Создаем схему reference вручную, так как Ent этого не делает автоматически
	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS reference"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed creating reference schema: %w", err)
	}

	// Создаем драйвер Ent на основе существующего соединения
	drv := entsql.OpenDB(dialect.Postgres, db)

	client := ent.NewClient(ent.Driver(drv))

	// Register global hooks
	client.Use(pg.SoftDeleteHook())

	return client, nil
}

// RunMigrations запускает автоматическую миграцию схем
func RunMigrations(client *ent.Client) error {
	ctx := context.Background()

	// При миграции Ent будет использовать SchemaConfig, заданный при инициализации клиента.
	// Опция WithDiffSchema(true) заставляет Ent учитывать схемы при сравнении текущего состояния БД и схемы Ent.
	if err := client.Schema.Create(ctx,
		schema.WithForeignKeys(true),
	); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	log.Println("ENT migrations completed successfully")
	return nil
}
