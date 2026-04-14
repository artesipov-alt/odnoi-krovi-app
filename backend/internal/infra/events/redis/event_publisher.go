package redisevent

import (
	"context"
	"encoding/json"
	"fmt"

	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	userevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
	"github.com/redis/go-redis/v9"
)

const channelBloodRequestCreated = "blood_request_created"
const channelDonorResponseApply = "donor_response_apply"
const channelDonationConfirmed = "donation_confirmed"
const channelRecipientResponseApply = "recipient_response_apply"
const channelDonorCancel = "donor_cancel"
const channelUserContact = "user_contact"

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

func (p *EventPublisher) PublishDonationConfirmed(
	ctx context.Context,
	event bloodsearchevent.DonationConfirmed,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelDonationConfirmed, payload).Err()
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

func (p *EventPublisher) PublishDonorCancel(
	ctx context.Context,
	event donorevent.DonorCancel,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelDonorCancel, payload).Err()
}

func (p *EventPublisher) PublishUserContact(
	ctx context.Context,
	event userevent.UserContact,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return p.client.Publish(ctx, channelUserContact, payload).Err()
}
