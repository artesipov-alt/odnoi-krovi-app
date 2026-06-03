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
		t.Fatal("TWIN24_BASE_URL is not set in .env file")
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

	authClient := twin24.NewClient(iamBaseURL)
	auth := twin24.NewAuth(authClient, *redisClient, os.Getenv("TWIN24_EMAIL"), os.Getenv("TWIN24_PASSWORD"))

	sender := twin24.NewOTPSender(baseURL, os.Getenv("TWIN24_BOT_SCENARIO_ID"), os.Getenv("TWIN24_BOT_CID"), authClient, auth)

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	// SendOTP создаст задание, добавит кандидата и запустит
	// Внимание: нужно заменить botScenarioID в коде на реальный
	err := sender.SendOTP(ctx, "79263165800", "2026")
	if err != nil {
		t.Fatalf("SendOTP failed: %v", err)
	}
	t.Log("SendOTP completed successfully")
}
