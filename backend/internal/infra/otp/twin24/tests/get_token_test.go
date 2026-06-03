package tests

import (
	"testing"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/otp/twin24"
	"github.com/redis/go-redis/v9"
)

func TestGetToken(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	if err := redisClient.Ping(t.Context()).Err(); err != nil {
		t.Fatalf("redisClient is not connected: %v", err)
	}

	client := twin24.NewClient("https://iam.twin24.ai/")
	auth := twin24.NewAuth(client, *redisClient, "Test02.06.2026@mail.ru", "f4f17fcd2715d4728dc4eca9674954e5")
	token, err := auth.GetToken(t.Context())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("expected token, got empty string")
	}
	t.Logf("token: %s", token)
}
