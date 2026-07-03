package events

import "time"

type BloodRequestCreated struct {
	RequestID  string    `json:"requestId"`
	BloodTypes []string  `json:"bloodTypes"`
	Regions    []string  `json:"regions"`
	TelegramID string    `json:"telegramId"`
	MaxID      string    `json:"maxId"`
	CreatedAt  time.Time `json:"createdAt"`
}
