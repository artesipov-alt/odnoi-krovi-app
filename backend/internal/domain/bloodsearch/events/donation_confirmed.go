package events

import "time"

type DonorInfo struct {
	UserName         string
	PetName          string
	ProviderMaxID    string
	ProviderTelegram string
	Phone            string
	BloodGroup       string
}

type RecipientInfo struct {
	PetName    string
	BloodGroup string
}

type DonationConfirmed struct {
	DonorData     DonorInfo
	RecipientData RecipientInfo
	Volume        float64
	CreatedAt     time.Time
}

func (e DonationConfirmed) EventName() string     { return "DonationConfirmed" }
func (e DonationConfirmed) OccurredAt() time.Time { return e.CreatedAt }
