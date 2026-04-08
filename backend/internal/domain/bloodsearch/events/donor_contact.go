package events

import (
	"time"
)

type DonorData struct {
	UserName         string
	PetName          string
	ProviderMaxID    string
	ProviderTelegram string
	Phone            string
	BloodGroup       string
}

type RecipientData struct {
	UserName         string
	PetName          string
	ProviderMaxID    string
	ProviderTelegram string
	Phone            string
	BloodGroup       string
	Volume           float64
}

type ApplyDonor struct {
	DonorData     DonorData
	RecipientData RecipientData
	CreatedAt     time.Time
}

func (e ApplyDonor) EventName() string     { return "ApplyDonor" }
func (e ApplyDonor) OccurredAt() time.Time { return e.CreatedAt }
