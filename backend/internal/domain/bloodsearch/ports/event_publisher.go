package ports

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
)

type EventPublisher interface {
	PublishBloodRequestCreated(ctx context.Context, event events.BloodRequestCreated) error
	PublishDonorApply(ctx context.Context, event events.ApplyDonor) error
}
