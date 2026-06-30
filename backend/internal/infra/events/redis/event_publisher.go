package redisevent

import (
	"context"
	"encoding/json"
	"fmt"

	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	userevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
	"github.com/redis/go-redis/v9"
)

const channelBloodRequestCreated = "blood_request_created"
const channelDonorResponseApply = "donor_response_apply"
const channelDonationConfirmed = "donation_confirmed"
const channelRecipientResponseApply = "recipient_response_apply"
const channelDonorCancel = "donor_cancel"
const channelDonorReject = "donor_reject"
const channelDonorNotConfirmed = "donor_not_confirmed"
const channelDonorCompleted = "donor_completed"
const channelUserContact = "user_contact"

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

func (p *EventPublisher) PublishBloodRequestCreated(
	ctx context.Context,
	event bloodsearchevent.BloodRequestCreated,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelBloodRequestCreated), payload).Err()
}

func (p *EventPublisher) PublishDonorApply(
	ctx context.Context,
	event bloodsearchevent.ApplyDonor,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelDonorResponseApply), payload).Err()
}

func (p *EventPublisher) PublishDonationConfirmed(
	ctx context.Context,
	event bloodsearchevent.DonationConfirmed,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelDonationConfirmed), payload).Err()
}

func (p *EventPublisher) PublishRecipientApply(
	ctx context.Context,
	event donorevent.RecipientApply,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelRecipientResponseApply), payload).Err()
}

func (p *EventPublisher) PublishDonorCancel(
	ctx context.Context,
	event donorevent.DonorCancel,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelDonorCancel), payload).Err()
}

func (p *EventPublisher) PublishDonorReject(
	ctx context.Context,
	event donorevent.DonorReject,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelDonorReject), payload).Err()
}

func (p *EventPublisher) PublishDonorNotConfirmed(
	ctx context.Context,
	event donorevent.DonorNotConfirmed,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelDonorNotConfirmed), payload).Err()
}

func (p *EventPublisher) PublishDonorCompleted(
	ctx context.Context,
	event donorevent.DonorCompleted,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelDonorCompleted), payload).Err()
}

func (p *EventPublisher) PublishUserContact(
	ctx context.Context,
	event userevent.UserContact,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, p.channel(channelUserContact), payload).Err()
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

func (p *NoOpEventPublisher) PublishBloodRequestCreated(ctx context.Context, event bloodsearchevent.BloodRequestCreated) error {
	return nil
}

func (p *NoOpEventPublisher) PublishDonorApply(ctx context.Context, event bloodsearchevent.ApplyDonor) error {
	return nil
}

func (p *NoOpEventPublisher) PublishDonationConfirmed(ctx context.Context, event bloodsearchevent.DonationConfirmed) error {
	return nil
}

func (p *NoOpEventPublisher) PublishRecipientApply(ctx context.Context, event donorevent.RecipientApply) error {
	return nil
}

func (p *NoOpEventPublisher) PublishDonorCancel(ctx context.Context, event donorevent.DonorCancel) error {
	return nil
}

func (p *NoOpEventPublisher) PublishDonorReject(ctx context.Context, event donorevent.DonorReject) error {
	return nil
}

func (p *NoOpEventPublisher) PublishDonorNotConfirmed(ctx context.Context, event donorevent.DonorNotConfirmed) error {
	return nil
}

func (p *NoOpEventPublisher) PublishDonorCompleted(ctx context.Context, event donorevent.DonorCompleted) error {
	return nil
}

func (p *NoOpEventPublisher) PublishUserContact(ctx context.Context, event userevent.UserContact) error {
	return nil
}

func (n *NoOpEventPublisher) PublishNotification(_ context.Context, _ ports.Notification) error {
	return nil
}
