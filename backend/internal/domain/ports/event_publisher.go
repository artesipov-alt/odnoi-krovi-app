package ports

import (
	"context"
	"time"

	bloodsearchevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	donorevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/events"
	userevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
)

type NotificationType string

const (
	NotifRecipientDonorWaiting    NotificationType = "recipient_donor_waiting"
	NotifRecipientInactiveWarning NotificationType = "recipient_inactive_warning"
	NotifRecipientSearchClosed    NotificationType = "recipient_search_closed"
	NotifRecipientEmptyShowcase   NotificationType = "recipient_empty_showcase"
	NotifDonorNotAccepted         NotificationType = "donor_not_accepted"
)

type Notification struct {
	Type      NotificationType `json:"type"`
	Targets   NotifTargets     `json:"targets"`
	Payload   map[string]any   `json:"payload"`
	CreatedAt time.Time        `json:"createdAt"`
}

type NotifTargets struct {
	TelegramID string `json:"telegramId,omitempty"`
	MaxID      string `json:"maxId,omitempty"`
}

type EventPublisher interface {
	PublishNotification(ctx context.Context, n Notification) error

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
	PublishDonorCompleted(ctx context.Context, event donorevent.DonorCompleted) error

	//==========================
	// 			User
	//==========================
	PublishUserContact(ctx context.Context, event userevent.UserContact) error
}
