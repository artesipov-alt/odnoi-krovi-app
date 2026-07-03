package events

import "time"

type RecipientApply struct {
	DonorName                       string    `json:"donorName"`
	DonorBloodGroup                 string    `json:"donorBloodGroup"`
	RecipientProviderMaxID          string    `json:"recipientProviderMaxId"`
	RecipientProviderTelegramID     string    `json:"recipientProviderTelegramId"`
	RecipientPetName                string    `json:"recipientPetName"`
	RecipientPetSearchingBloodGroup []string  `json:"recipientPetSearchingBloodGroup"`
	RecipientPetNeededVolume        float64   `json:"recipientPetNeededVolume"`
	CreatedAt                       time.Time `json:"createdAt"`
}
