package config

import (
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig содержит настройки подключения к базе данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabaseConfig создает конфигурацию базы данных из переменных окружения
func NewDatabaseConfig() *DatabaseConfig {
	port, _ := strconv.Atoi(GetEnv("DB_PORT", "5432"))

	return &DatabaseConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", "postgres"),
		DBName:   GetEnv("DB_NAME", "odnoi_krovi"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}
}

// GetDSN возвращает строку подключения к PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// ConnectDB устанавливает соединение с базой данных
func ConnectDB() (*gorm.DB, error) {
	config := NewDatabaseConfig()

	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
