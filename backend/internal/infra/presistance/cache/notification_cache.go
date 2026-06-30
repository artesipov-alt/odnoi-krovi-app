package cache

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
)

type NotificationCache interface {
	WasSent(ctx context.Context, notifType ports.NotificationType, entityID string) (bool, error)
	MarkSent(ctx context.Context, notifType ports.NotificationType, entityID string, ttl time.Duration) error
	MarkYesPressed(ctx context.Context, bloodRequestID string) error
	WasYesPressed(ctx context.Context, bloodRequestID string) (bool, error)
}
