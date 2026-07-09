package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/redis/go-redis/v9"
)

type NotificationCache struct {
	client *redis.Client
}

func NewNotificationCache(client *redis.Client) *NotificationCache {
	return &NotificationCache{
		client: client,
	}
}

func (c *NotificationCache) WasSent(ctx context.Context, notifType ports.NotificationType, entityID string) (bool, error) {
	key := fmt.Sprintf("notif:sent:%s:%s", notifType, entityID)
	err := c.client.Get(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("WasSent: %w", err)
	}
	return true, nil
}

func (c *NotificationCache) MarkSent(ctx context.Context, notifType ports.NotificationType, entityID string, ttl time.Duration) error {
	key := fmt.Sprintf("notif:sent:%s:%s", notifType, entityID)
	if err := c.client.Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf("MarkSent: %w", err)
	}
	return nil
}

func (c *NotificationCache) MarkYesPressed(ctx context.Context, bloodRequestID string) error {
	key := fmt.Sprintf("notif:yes_pressed:%s", bloodRequestID)
	if err := c.client.Set(ctx, key, "1", 48*time.Hour).Err(); err != nil {
		return fmt.Errorf("MarkYesPressed: %w", err)
	}
	return nil
}

func (c *NotificationCache) WasYesPressed(ctx context.Context, bloodRequestID string) (bool, error) {
	key := fmt.Sprintf("notif:yes_pressed:%s", bloodRequestID)
	err := c.client.Get(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("WasYesPressed: %w", err)
	}
	return true, nil
}
