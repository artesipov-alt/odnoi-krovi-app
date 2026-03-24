package events

import "time"

type ApplyDonor struct {
	DonorData     EventData
	RecipientData EventData
	CreatedAt     time.Time
}

type EventData struct {
	Name             string
	ProviderMaxID    string
	ProviderTelegram string
	Phone            string
}

func (e ApplyDonor) EventName() string     { return "ApplyDonor" }
func (e ApplyDonor) OccurredAt() time.Time { return e.CreatedAt }
