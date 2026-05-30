package twin24

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenCacheKey = "twin24:token"

type Auth struct {
	client   *Client
	redis    redis.Client
	email    string
	password string
}

func NewAuth(client *Client, redis redis.Client, email, password string) *Auth {
	return &Auth{
		client:   client,
		redis:    redis,
		email:    email,
		password: password,
	}
}

func (a *Auth) GetToken(ctx context.Context) (string, error) {
	// достаём из кеша
	token, err := a.redis.Get(ctx, tokenCacheKey).Result()
	if err == nil {
		return token, nil
	}

	// кеша нет — логинимся
	return a.login(ctx)
}

func (a *Auth) RefreshToken(ctx context.Context) (string, error) {
	// сбрасываем кеш и логинимся заново
	a.redis.Del(ctx, tokenCacheKey)
	return a.login(ctx)
}

func (a *Auth) login(ctx context.Context) (string, error) {
	resp, err := a.client.Do(ctx, http.MethodPost, "/api/v1/auth/login", map[string]any{
		"email":    a.email,
		"password": a.password,
		"ttl":      36000, // 10 часов
	})
	if err != nil {
		return "", fmt.Errorf("twin24 login: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("twin24 login decode: %w", err)
	}

	// кешируем с TTL чуть меньше чем у токена
	a.redis.Set(ctx, tokenCacheKey, result.Token, 9*time.Hour)

	return result.Token, nil
}
