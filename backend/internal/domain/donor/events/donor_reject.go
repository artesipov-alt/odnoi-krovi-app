package events

import "time"

type DonorReject struct {
	RecipientPetName        string
	RecipientBloodGroup     string
	DonorProviderMaxID      string
	DonorProviderTelegramID string
	DonorPetName            string
	RejectedReason          string
	CreatedAt               time.Time
}

func (e DonorReject) EventName() string     { return "DonorReject" }
func (e DonorReject) OccurredAt() time.Time { return e.CreatedAt }
