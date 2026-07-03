package events

import "time"

type DonorInfo struct {
	UserName         string `json:"userName"`
	PetName          string `json:"petName"`
	ProviderMaxID    string `json:"providerMaxId"`
	ProviderTelegram string `json:"providerTelegram"`
	Phone            string `json:"phone"`
	BloodGroup       string `json:"bloodGroup"`
}

type RecipientInfo struct {
	PetName    string `json:"petName"`
	BloodGroup string `json:"bloodGroup"`
}

type DonationConfirmed struct {
	DonorData     DonorInfo     `json:"donorData"`
	RecipientData RecipientInfo `json:"recipientData"`
	Volume        float64       `json:"volume"`
	CreatedAt     time.Time     `json:"createdAt"`
}
