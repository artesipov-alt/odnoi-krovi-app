package redisevent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/redis/go-redis/v9"
)

const channelBloodRequestCreated = "blood_request.created"

type EventPublisher struct {
	client *redis.Client
}

func NewEventPublisher(client *redis.Client) *EventPublisher {
	return &EventPublisher{client: client}
}

func (p *EventPublisher) PublishBloodRequestCreated(
	ctx context.Context,
	event events.BloodRequestCreated,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelBloodRequestCreated, payload).Err()
}
