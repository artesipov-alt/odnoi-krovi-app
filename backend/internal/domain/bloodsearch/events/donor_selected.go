package events

import "time"

// DonorSelected — реципиент выбрал донора из списка потенциальных.
// Уведомление отправляется только донору.
type DonorSelected struct {
	DonorData     DonorSelectedDonorData     `json:"donorData"`
	RecipientData DonorSelectedRecipientData `json:"recipientData"`
	CreatedAt     time.Time                  `json:"createdAt"`
}

// DonorSelectedDonorData — данные донора для роутинга уведомления.
type DonorSelectedDonorData struct {
	ProviderMaxID    string `json:"providerMaxId"`
	ProviderTelegram string `json:"providerTelegram"`
}

// DonorSelectedRecipientData — данные реципиента для текста уведомления.
type DonorSelectedRecipientData struct {
	PetName    string  `json:"petName"`
	PetType    string  `json:"petType"` // "cat" или "dog"
	Volume     float64 `json:"volume"`  // мл
	BloodGroup string  `json:"bloodGroup"`
}
