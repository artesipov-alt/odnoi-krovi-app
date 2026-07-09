package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/otp/twin24"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/redis/otp"

	"crypto/rand"
	"math/big"
)

type ChangePhoneHandler struct {
	userRepo  user.Repository
	otpRepo   otp.OTPRepository
	otpSender *twin24.OTPSender
}

func NewChangePhoneHandler(userRepo user.Repository, otpRepo otp.OTPRepository, otpSender *twin24.OTPSender) *ChangePhoneHandler {
	return &ChangePhoneHandler{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		otpSender: otpSender,
	}
}

func (h *ChangePhoneHandler) Handle(ctx context.Context, userID string, phone string) error {
	otpCode, err := generateOTPCode()
	if err != nil {
		return fmt.Errorf("change phone: %w", err)
	}

	if err := h.otpRepo.Save(ctx, userID, otp.NewOTPData(otpCode, phone), 5*time.Minute); err != nil {
		return fmt.Errorf("change phone: %w", err)
	}

	// отправляем OTP через Twin24
	if err := h.otpSender.SendOTP(ctx, phone, otpCode); err != nil {
		// не удаляем OTP из Redis, чтобы можно было попробовать снова
		return fmt.Errorf("change phone: send otp: %w", err)
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
