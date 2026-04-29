package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/migrate"
	_ "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/runtime"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
	sloghttp "github.com/samber/slog-http"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	schemaent "entgo.io/ent/dialect/sql/schema"
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
func NewEntConfig(env string) *EntConfig {
	var dbname string
	switch env {
	case "PROD", "prod", "production":
		dbname = os.Getenv("DB_NAME")
	default:
		dbname = os.Getenv("DB_NAME_DEV")
	}

	if dbname == "" {
		slog.Error("Не указано название базы данных в окружении env")
		os.Exit(1)
	}

	return &EntConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   dbname,
		SSLMode:  os.Getenv("DB_SSLMODE"),
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

// GetDSN возвращает строку подключения для PostgreSQL
func (c *EntConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func ConnectEnt(config *EntConfig) (*ent.Client, *sql.DB, error) {
	dsn := config.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	debugDrv := dialect.DebugWithContext(drv, func(ctx context.Context, v ...any) {
		slog.DebugContext(ctx, "SQL",
			"query", fmt.Sprint(v...),
			"rid", sloghttp.GetRequestIDFromContext(ctx),
		)
	})

	client := ent.NewClient(ent.Driver(debugDrv))
	client.Intercept(schema.DbInterceptor())
	client.Use(schema.SoftDeleteHook())

	return client, db, nil
}

// RunMigrations запускает автоматическую миграцию схем
func RunMigrations(client *ent.Client, db *sql.DB) error {
	ctx := context.Background()

	// При миграции Ent будет использовать SchemaConfig, заданный при инициализации клиента.
	// Опция WithDiffSchema(true) заставляет Ent учитывать схемы при сравнении текущего состояния БД и схемы Ent.
	if err := client.Schema.Create(ctx,
		schemaent.WithForeignKeys(true),
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	log.Println("ENT migrations completed successfully")

	// Создаём partial unique index для поддержки множественных closed/draft заявок на одного питомца
	// Standard unique constraint был удалён WithDropIndex, создаём partial только для active заявок
	if _, err := db.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS blood_requests_pet_id_key
		ON blood_requests (pet_id) WHERE status = 'active'
	`); err != nil {
		return fmt.Errorf("failed to create partial unique index: %w", err)
	}
	log.Println("Custom partial unique index created successfully")
	return nil
}
