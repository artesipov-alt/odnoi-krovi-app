package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config содержит конфигурацию для Redis
type Config struct {
	// RedisAddr адрес Redis сервера (host:port)
	RedisAddr string
	// RedisPassword пароль для Redis
	RedisPassword string
	// RedisDB номер базы данных Redis
	RedisDB int
	// RedisMaxRetries максимальное количество попыток переподключения
	RedisMaxRetries int
	// RedisPoolSize размер пула соединений
	RedisPoolSize int
}

// NewRedisClient создает новый Redis клиент
func NewRedisClient(addr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}

// NewRedisClientFromEnv создает Redis клиент из переменных окружения
func NewRedisClientFromEnv() (*redis.Client, error) {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisAddr := getEnv("REDIS_ADDR", redisHost+":"+redisPort)

	config := Config{
		RedisAddr:     redisAddr,
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
	}

	return NewRedisClient(config.RedisAddr, config.RedisPassword, config.RedisDB)
}

// NewDefaultRedisClient создает Redis клиент с настройками по умолчанию
func NewDefaultRedisClient() (*redis.Client, error) {
	config := Config{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
	}

	return NewRedisClient(config.RedisAddr, config.RedisPassword, config.RedisDB)
}

// getEnv получает переменную окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как целое число или значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
