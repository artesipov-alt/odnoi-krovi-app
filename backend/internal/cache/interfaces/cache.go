package interfaces

import (
	"context"
	"errors"
	"time"
)

// ErrCacheMiss возвращается когда ключ не найден в кэше
var ErrCacheMiss = errors.New("cache miss")

// Cache определяет интерфейс для операций с кэшем
type Cache interface {
	// Get получает значение по ключу
	Get(ctx context.Context, key string) ([]byte, error)

	// Set сохраняет значение по ключу с TTL
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete удаляет значение по ключу
	Delete(ctx context.Context, key string) error

	// Exists проверяет существование ключа
	Exists(ctx context.Context, key string) (bool, error)

	// Expire устанавливает TTL для ключа
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// GetTTL получает оставшееся время жизни ключа
	GetTTL(ctx context.Context, key string) (time.Duration, error)

	// Increment увеличивает числовое значение по ключу
	Increment(ctx context.Context, key string, value int64) (int64, error)

	// Decrement уменьшает числовое значение по ключу
	Decrement(ctx context.Context, key string, value int64) (int64, error)

	// Keys возвращает список ключей по паттерну
	Keys(ctx context.Context, pattern string) ([]string, error)

	// Flush очищает все ключи
	Flush(ctx context.Context) error

	// Ping проверяет соединение с кэшем
	Ping(ctx context.Context) error

	// Close закрывает соединение с кэшем
	Close() error

	// GetJSON получает JSON объект по ключу и десериализует его
	GetJSON(ctx context.Context, key string, target interface{}) error

	// SetJSON сериализует объект в JSON и сохраняет его по ключу с TTL
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}
