package events

import "time"

type ContactData struct {
	Name             string `json:"name"`
	ProviderMaxID    string `json:"providerMaxId"`
	ProviderTelegram string `json:"providerTelegram"`
	Phone            string `json:"phone"`
}

type DonorNotConfirmed struct {
	DonorPetName        string      `json:"donorPetName"`
	DonorBloodGroup     string      `json:"donorBloodGroup"`
	DonorProviderMaxID  string      `json:"donorProviderMaxId"`
	RecipientPetName    string      `json:"recipientPetName"`
	RecipientBloodGroup string      `json:"recipientBloodGroup"`
	RecipientUserData   ContactData `json:"recipientUserData"`
	DonorUserData       ContactData `json:"donorUserData"`
	CreatedAt           time.Time   `json:"createdAt"`
}
