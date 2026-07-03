package events

import "time"

type DonorCompleted struct {
	DonorPetName                string    `json:"donorPetName"`
	DonorBloodGroup             string    `json:"donorBloodGroup"`
	RecipientProviderMaxID      string    `json:"recipientProviderMaxId"`
	RecipientProviderTelegramID string    `json:"recipientProviderTelegramId"`
	RecipientPetName            string    `json:"recipientPetName"`
	Amount                      float64   `json:"amount"`
	CreatedAt                   time.Time `json:"createdAt"`
}
