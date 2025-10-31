package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/cache/interfaces"
	"github.com/redis/go-redis/v9"
)

// RedisCache реализует интерфейс Cache с использованием Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache создает новый экземпляр RedisCache
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
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

	return &RedisCache{
		client: client,
	}, nil
}

// Get получает значение по ключу
func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, interfaces.ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return result, nil
}

// Set сохраняет значение по ключу с TTL
func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Delete удаляет значение по ключу
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists проверяет существование ключа
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return result > 0, nil
}

// Expire устанавливает TTL для ключа
func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set expire for key %s: %w", key, err)
	}
	return nil
}

// GetTTL получает оставшееся время жизни ключа
func (r *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}
	return result, nil
}

// Increment увеличивает числовое значение по ключу
func (r *RedisCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	result, err := r.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return result, nil
}

// Decrement уменьшает числовое значение по ключу
func (r *RedisCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	result, err := r.client.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s: %w", key, err)
	}
	return result, nil
}

// Keys возвращает список ключей по паттерну
func (r *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	result, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}
	return result, nil
}

// Flush очищает все ключи
func (r *RedisCache) Flush(ctx context.Context) error {
	if err := r.client.FlushDB(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}

// Ping проверяет соединение с кэшем
func (r *RedisCache) Ping(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

// Close закрывает соединение с кэшем
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GetJSON получает JSON объект по ключу и десериализует его
func (r *RedisCache) GetJSON(ctx context.Context, key string, target interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON for key %s: %w", key, err)
	}
	return nil
}

// SetJSON сериализует объект в JSON и сохраняет его по ключу с TTL
func (r *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for key %s: %w", key, err)
	}

	return r.Set(ctx, key, data, ttl)
}
