package events

import "time"

type DonorCompleted struct {
	DonorPetName        string
	DonorBloodGroup     string
	RecipientProviderMaxID string
	RecipientPetName    string
	Amount              float64
	CreatedAt           time.Time
}

func (e DonorCompleted) EventName() string     { return "DonorCompleted" }
func (e DonorCompleted) OccurredAt() time.Time { return e.CreatedAt }
