package events

import "time"

type RecipientApply struct {
	DonorName       string
	DonorBloodGroup string
	CreatedAt       time.Time
}

func (e RecipientApply) EventName() string     { return "RecipientApply" }
func (e RecipientApply) OccurredAt() time.Time { return e.CreatedAt }
