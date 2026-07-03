package events

import "time"

type DonorCancel struct {
	DonorName                   string    `json:"donorName"`
	DonorBloodGroup             string    `json:"donorBloodGroup"`
	RecipientProviderMaxID      string    `json:"recipientProviderMaxId"`
	RecipientProviderTelegramID string    `json:"recipientProviderTelegramId"`
	RecipientPetName            string    `json:"recipientPetName"`
	CreatedAt                   time.Time `json:"createdAt"`
}
