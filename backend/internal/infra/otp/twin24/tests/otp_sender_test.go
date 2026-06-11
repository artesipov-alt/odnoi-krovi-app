package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/otp/twin24"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func TestOTPSender_SendOTP(t *testing.T) {
	// загружаем .env из корня проекта (на уровень выше backend/)
	// от backend/internal/infra/otp/twin24/tests/ нужно подняться на 6 уровней
	if err := godotenv.Load("../../../../../../.env"); err != nil {
		t.Logf("Warning: .env not loaded: %v", err)
	}

	baseURL := os.Getenv("TWIN24_BASE_URL")
	if baseURL == "" {
		baseURL = "https://twin24.ai"
	}

	iamBaseURL := os.Getenv("TWIN24_IAM_URL")
	if iamBaseURL == "" {
		iamBaseURL = "https://iam.twin24.ai"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	if err := redisClient.Ping(t.Context()).Err(); err != nil {
		t.Fatalf("redisClient is not connected: %v", err)
	}

	// Создаем OTP Sender используя переменные окружения
	sender := twin24.NewOTPSenderFromEnv(redisClient)
	if sender == nil {
		t.Fatal("OTPSender is nil - check TWIN24_* env vars")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	// SendOTP создаст задание, добавит кандидата и запустит
	// Внимание: нужно заменить botScenarioID в коде на реальный
	err := sender.SendOTP(ctx, "+79264187658", "1111")
	if err != nil {
		t.Fatalf("SendOTP failed: %v", err)
	}
	t.Log("SendOTP completed successfully")
}
