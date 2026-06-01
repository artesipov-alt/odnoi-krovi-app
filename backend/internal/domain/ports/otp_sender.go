package ports

import "context"

type OTPSender interface {
	SendOTP(ctx context.Context, phone string, code string) error
}
