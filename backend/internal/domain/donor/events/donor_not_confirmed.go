package events

import "time"

type DonorNotConfirmed struct {
	DonorPetName        string
	DonorBloodGroup     string
	DonorProviderMaxID  string
	RecipientPetName    string
	RecipientBloodGroup string
	RecipientUserData   ContactData
	DonorUserData       ContactData
	CreatedAt           time.Time
}

type ContactData struct {
	Name             string
	ProviderMaxID    string
	ProviderTelegram string
	Phone            string
}

func (e DonorNotConfirmed) EventName() string     { return "DonorNotConfirmed" }
func (e DonorNotConfirmed) OccurredAt() time.Time { return e.CreatedAt }
