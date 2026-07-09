package redisevent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/redis/go-redis/v9"
)

const channelEvents = "events"
const channelNotifications = "notifications"

type EventPublisher struct {
	client    *redis.Client
	envPrefix string
}

func NewEventPublisher(client *redis.Client, env string) *EventPublisher {
	prefix := ""
	if env == "development" || env == "dev" {
		prefix = "dev:"
	} else {
		prefix = "prod:"
	}

	return &EventPublisher{
		client:    client,
		envPrefix: prefix,
	}
}

func (p *EventPublisher) channel(name string) string {
	return p.envPrefix + name
}

func (p *EventPublisher) PublishEvent(ctx context.Context, eventType ports.EventType, payload any) error {
	envelope := ports.EventEnvelope{
		Type:      eventType,
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("PublishEvent: marshal: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelEvents), data).Err()
}

func (p *EventPublisher) PublishNotification(ctx context.Context, n ports.Notification) error {
	payload, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("PublishNotification: marshal: %w", err)
	}

	if err := p.client.Publish(ctx, p.channel(channelNotifications), payload).Err(); err != nil {
		return fmt.Errorf("PublishNotification: publish: %w", err)
	}

	return nil
}

// NoOpEventPublisher is a no-operation event publisher that does nothing.
// Used when Redis is not available, allowing the server to start without event publishing.
type NoOpEventPublisher struct{}

func (n *NoOpEventPublisher) PublishEvent(_ context.Context, _ ports.EventType, _ any) error {
	return nil
}

func (n *NoOpEventPublisher) PublishNotification(_ context.Context, _ ports.Notification) error {
	return nil
}
