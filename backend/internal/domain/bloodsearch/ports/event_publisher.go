package ports

import (
	"context"

	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
)

type EventPublisher interface {
	PublishBloodRequestCreated(ctx context.Context, event bloodsearchevent.BloodRequestCreated) error
	PublishDonorApply(ctx context.Context, event bloodsearchevent.ApplyDonor) error
	PublishRecipientApply(ctx context.Context, event donorevent.RecipientApply) error
}
