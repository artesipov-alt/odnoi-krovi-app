package events

import "time"

type DonorData struct {
	UserName         string `json:"userName"`
	PetName          string `json:"petName"`
	ProviderMaxID    string `json:"providerMaxId"`
	ProviderTelegram string `json:"providerTelegram"`
	Phone            string `json:"phone"`
	BloodGroup       string `json:"bloodGroup"`
}

type RecipientData struct {
	UserName         string  `json:"userName"`
	PetName          string  `json:"petName"`
	ProviderMaxID    string  `json:"providerMaxId"`
	ProviderTelegram string  `json:"providerTelegram"`
	Phone            string  `json:"phone"`
	BloodGroup       string  `json:"bloodGroup"`
	Volume           float64 `json:"volume"`
}

type ApplyDonor struct {
	DonorData     DonorData     `json:"donorData"`
	RecipientData RecipientData `json:"recipientData"`
	CreatedAt     time.Time     `json:"createdAt"`
}
