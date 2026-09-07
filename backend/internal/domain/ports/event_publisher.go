package ports

import (
	"context"
	"time"
)

type NotificationType string

const (
	NotifRecipientDonorWaiting         NotificationType = "recipient_donor_waiting"
	NotifRecipientInactiveWarning      NotificationType = "recipient_inactive_warning"
	NotifRecipientSearchClosed         NotificationType = "recipient_search_closed"
	NotifRecipientEmptyShowcase        NotificationType = "recipient_empty_showcase"
	NotifDonorNotAccepted              NotificationType = "donor_not_accepted"
	NotifRecipientSearchClosedInactive NotificationType = "recipient_search_closed_inactive" // п.6
	NotifRecipientAcceptedReminder12h  NotificationType = "recipient_accepted_reminder_12h"
	NotifDonorAcceptedReminder12h      NotificationType = "donor_accepted_reminder_12h"
	NotifRecipientAcceptedReminder24h  NotificationType = "recipient_accepted_reminder_24h"
	NotifUserVerifiedNoPets            NotificationType = "user_verified_no_pets"
	NotifUserNotVerified               NotificationType = "user_not_verified"
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

type EventType string

const (
	EventBloodRequestCreated EventType = "blood_request_created"
	EventDonorApply          EventType = "donor_response_apply"
	EventDonorSelected       EventType = "donor_selected"
	EventDonationConfirmed   EventType = "donation_confirmed"
	EventRecipientApply      EventType = "recipient_response_apply"
	EventDonorCancel         EventType = "donor_cancel"
	EventDonorReject         EventType = "donor_reject"
	EventDonorNotConfirmed   EventType = "donor_not_confirmed"
	EventDonorCompleted      EventType = "donor_completed"
	EventUserContact         EventType = "user_contact"
)

type EventEnvelope struct {
	Type      EventType `json:"type"`
	Payload   any       `json:"payload"`
	CreatedAt time.Time `json:"createdAt"`
}

type EventPublisher interface {
	PublishEvent(ctx context.Context, eventType EventType, payload any) error
	PublishNotification(ctx context.Context, n Notification) error
}
