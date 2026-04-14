package ports

import (
	"context"

	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	userevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
)

type EventPublisher interface {
	//==========================
	// 			BloodSearch
	//==========================
	PublishBloodRequestCreated(ctx context.Context, event bloodsearchevent.BloodRequestCreated) error
	PublishDonorApply(ctx context.Context, event bloodsearchevent.ApplyDonor) error
	PublishDonationConfirmed(ctx context.Context, event bloodsearchevent.DonationConfirmed) error

	//==========================
	// 			Donor
	//==========================
	PublishRecipientApply(ctx context.Context, event donorevent.RecipientApply) error
	PublishDonorCancel(ctx context.Context, event donorevent.DonorCancel) error
	PublishDonorReject(ctx context.Context, event donorevent.DonorReject) error
	PublishDonorNotConfirmed(ctx context.Context, event donorevent.DonorNotConfirmed) error

	//==========================
	// 			User
	//==========================
	PublishUserContact(ctx context.Context, event userevent.UserContact) error
}
