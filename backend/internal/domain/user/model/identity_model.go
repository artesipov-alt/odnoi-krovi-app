package model

import (
	"time"
)

// Identity представляет доменную модель пользователя
type Identity struct {
	ID             string
	UserID         string
	ProviderName   string // Assuming useridentity.Provider can be represented as a string
	ProviderUserID int64
	XBToken        string
	Metadata       map[string]any
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// NewUser creates a new User aggregate with validation
func NewIdentity(providerID int64, providerName, authBotToken string) (*Identity, error) {

	idn := &Identity{
		ProviderUserID: providerID,
		ProviderName:   providerName,
		XBToken:        authBotToken,
	}

	return idn, nil
}
