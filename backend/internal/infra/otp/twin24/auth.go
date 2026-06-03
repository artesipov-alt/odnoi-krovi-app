package twin24

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	tokenCacheKey        = "twin24:token"
	refreshTokenCacheKey = "twin24:refreshToken"
)

type loginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

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
	token, err := a.redis.Get(ctx, tokenCacheKey).Result()
	if err == nil && token != "" {
		return token, nil
	}
	// удаляем пустой/невалидный кеш
	if err == nil && token == "" {
		a.redis.Del(ctx, tokenCacheKey)
	}
	return a.login(ctx)
}

func (a *Auth) RefreshToken(ctx context.Context) (string, error) {
	token, err := a.refreshViaAPI(ctx)
	if err == nil {
		return token, nil
	}
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("twin24 login: unexpected status %d, body: %s", resp.StatusCode, string(body))
	}

	var result loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("twin24 login decode: %w", err)
	}

	if result.Token == "" {
		return "", fmt.Errorf("twin24 login: empty token in response")
	}

	a.redis.Set(ctx, tokenCacheKey, result.Token, 9*time.Hour)
	if result.RefreshToken != "" {
		a.redis.Set(ctx, refreshTokenCacheKey, result.RefreshToken, 340*24*time.Hour)
	}

	return result.Token, nil
}

type refreshResponse struct {
	Token string `json:"token"`
}

func (a *Auth) refreshViaAPI(ctx context.Context) (string, error) {
	refreshToken, err := a.redis.Get(ctx, refreshTokenCacheKey).Result()
	if err != nil {
		return "", fmt.Errorf("refreshToken not found in cache: %w", err)
	}

	resp, err := a.client.Do(ctx, http.MethodPost, "/api/v1/auth/refresh", map[string]any{
		"refreshToken": refreshToken,
		"ttl":          3600,
	})
	if err != nil {
		return "", fmt.Errorf("twin24 refresh: %w", err)
	}
	defer resp.Body.Close()

	// 401 — refresh token невалиден, нужен перелогин
	if resp.StatusCode == http.StatusUnauthorized {
		a.redis.Del(ctx, refreshTokenCacheKey)
		return "", fmt.Errorf("refresh token invalid (401)")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("twin24 refresh: unexpected status %d, body: %s", resp.StatusCode, string(body))
	}

	var result refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("twin24 refresh decode: %w", err)
	}

	if result.Token == "" {
		return "", fmt.Errorf("twin24 refresh: empty token in response")
	}

	// TTL=3600 сек = 1 час
	a.redis.Set(ctx, tokenCacheKey, result.Token, 1*time.Hour)

	return result.Token, nil
}
