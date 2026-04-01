package events

import (
	"time"

	userevents "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
)

type ApplyDonor struct {
	DonorData     userevents.ContactData
	RecipientData userevents.ContactData
	CreatedAt     time.Time
}

func (e ApplyDonor) EventName() string     { return "ApplyDonor" }
func (e ApplyDonor) OccurredAt() time.Time { return e.CreatedAt }
