package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/redis"

	"crypto/rand"
	"math/big"
)

type ChangePhoneHandler struct {
	userRepo user.Repository
	otpRepo  redis.OTPRepository
}

func NewChangePhoneHandler(userRepo user.Repository, otpRepo redis.OTPRepository) *ChangePhoneHandler {
	return &ChangePhoneHandler{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (h *ChangePhoneHandler) Handle(ctx context.Context, userID string, phone string) error {
	// otpCode, err := generateOTPCode()
	// if err != nil {
	// 	return fmt.Errorf("change phone: %w", err)
	// }

	otpData := redis.OTPData{
		NewPhone: phone,
		Code:     "2026",
	}
	if err := h.otpRepo.Save(ctx, userID, otpData, 5*time.Minute); err != nil {
		return fmt.Errorf("change phone: %w", err)
	}

	return nil
}

func generateOTPCode() (string, error) {
	const length = 4
	max := big.NewInt(10000) // 0000 - 9999

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}

	// форматируем с ведущими нулями — "0042" а не "42"
	return fmt.Sprintf("%04d", n), nil
}
