package events

import "time"

type DonorReject struct {
	RecipientPetName        string    `json:"recipientPetName"`
	RecipientBloodGroup     string    `json:"recipientBloodGroup"`
	DonorProviderMaxID      string    `json:"donorProviderMaxId"`
	DonorProviderTelegramID string    `json:"donorProviderTelegramId"`
	DonorPetName            string    `json:"donorPetName"`
	RejectedReason          string    `json:"rejectedReason"`
	CreatedAt               time.Time `json:"createdAt"`
}
