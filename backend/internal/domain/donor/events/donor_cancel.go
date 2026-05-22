package events

import "time"

type DonorCancel struct {
	DonorName                   string
	DonorBloodGroup             string
	RecipientProviderMaxID      string
	RecipientProviderTelegramID string
	RecipientPetName            string
	CreatedAt                   time.Time
}

func (e DonorCancel) EventName() string     { return "DonorCancel" }
func (e DonorCancel) OccurredAt() time.Time { return e.CreatedAt }
