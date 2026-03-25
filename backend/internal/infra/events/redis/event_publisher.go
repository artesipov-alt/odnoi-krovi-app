package redisevent

import (
	"context"
	"encoding/json"
	"fmt"

	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	"github.com/redis/go-redis/v9"
)

const channelBloodRequestCreated = "blood_request.created"
const channelDonorResponseApply = "donor_response.apply"
const channelRecipientResponseApply = "recipient_response.apply"

type EventPublisher struct {
	client *redis.Client
}

func NewEventPublisher(client *redis.Client) *EventPublisher {
	return &EventPublisher{client: client}
}

func (p *EventPublisher) PublishBloodRequestCreated(
	ctx context.Context,
	event bloodsearchevent.BloodRequestCreated,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelBloodRequestCreated, payload).Err()
}

func (p *EventPublisher) PublishDonorApply(
	ctx context.Context,
	event bloodsearchevent.ApplyDonor,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelDonorResponseApply, payload).Err()
}

func (p *EventPublisher) PublishRecipientApply(
	ctx context.Context,
	event donorevent.RecipientApply,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelRecipientResponseApply, payload).Err()
}
