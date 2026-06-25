package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPData struct {
	Code     string `json:"code"`
	NewPhone string `json:"newPhone"`
}

func NewOTPData(code, newPhone string) OTPData {
	return OTPData{
		Code:     code,
		NewPhone: newPhone,
	}
}

type OTPRepository interface {
	Save(ctx context.Context, userID string, data OTPData, ttl time.Duration) error
	Get(ctx context.Context, userID string) (OTPData, error)
	Delete(ctx context.Context, userID string) error
}

type OTPRepo struct {
	client *redis.Client
}

func NewOTPRepo(client *redis.Client) *OTPRepo {
	return &OTPRepo{
		client: client,
	}
}

func (r *OTPRepo) Save(ctx context.Context, userID string, data OTPData, ttl time.Duration) error {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, userID, marshalledData, ttl).Err()
}

func (r *OTPRepo) Get(ctx context.Context, userID string) (OTPData, error) {
	var otpData OTPData
	val, err := r.client.Get(ctx, userID).Result()
	if err != nil {
		return otpData, err
	}
	err = json.Unmarshal([]byte(val), &otpData)
	return otpData, err
}

func (r *OTPRepo) Delete(ctx context.Context, userID string) error {
	return r.client.Del(ctx, userID).Err()
}

// NoOpOTPRepo — заглушка для случаев, когда Redis недоступен.
// Все операции возвращают ошибку, сигнализируя, что OTP-сервис отключён.
type NoOpOTPRepo struct{}

func NewNoOpOTPRepo() *NoOpOTPRepo {
	return &NoOpOTPRepo{}
}

func (r *NoOpOTPRepo) Save(_ context.Context, _ string, _ OTPData, _ time.Duration) error {
	return fmt.Errorf("OTP service is unavailable: Redis is not connected")
}

func (r *NoOpOTPRepo) Get(_ context.Context, _ string) (OTPData, error) {
	return OTPData{}, fmt.Errorf("OTP service is unavailable: Redis is not connected")
}

func (r *NoOpOTPRepo) Delete(_ context.Context, _ string) error {
	return fmt.Errorf("OTP service is unavailable: Redis is not connected")
}
