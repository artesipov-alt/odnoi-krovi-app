package redis

import (
	"fmt"
	"os"
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/cache/interfaces"
)

// Config содержит конфигурацию для кэша
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

// NewCache создает новый экземпляр кэша на основе конфигурации
func NewCache(config Config) (interfaces.Cache, error) {
	// Валидация конфигурации
	if config.RedisAddr == "" {
		return nil, fmt.Errorf("redis address is required")
	}

	// Создаем Redis кэш
	cache, err := NewRedisCache(
		config.RedisAddr,
		config.RedisPassword,
		config.RedisDB,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis cache: %w", err)
	}

	return cache, nil
}

// NewCacheFromEnv создает кэш из переменных окружения
func NewCacheFromEnv() (interfaces.Cache, error) {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisAddr := getEnv("REDIS_ADDR", redisHost+":"+redisPort)

	config := Config{
		RedisAddr:       redisAddr,
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
		RedisMaxRetries: getEnvAsInt("REDIS_MAX_RETRIES", 3),
		RedisPoolSize:   getEnvAsInt("REDIS_POOL_SIZE", 10),
	}

	return NewCache(config)
}

// NewDefaultCache создает кэш с настройками по умолчанию
func NewDefaultCache() (interfaces.Cache, error) {
	config := Config{
		RedisAddr:       "localhost:6379",
		RedisPassword:   "",
		RedisDB:         0,
		RedisMaxRetries: 3,
		RedisPoolSize:   10,
	}

	return NewCache(config)
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
