package bonus

import "time"

// Bonus represents a bonus entity in the domain.
type Bonus struct {
	ID           string
	UserID       *string
	PartnerName  string
	Description  string
	Target       string
	Recipient    string
	Category     string
	Subcategory  *string
	PromoCode    string
	ExpiresAt    time.Time
	PlatformName string
	PlatformURL  *string
	Stage        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// NewLockBonus creates a new lock bonus for users who donated recently
func NewLockBonus() *Bonus {
	return &Bonus{
		PartnerName: "Портал",
		Description: "Пользователь уже получал свои бонусы в течение двух месяцев.",
		Category:    "lock",
		Target:      "all",
		Recipient:   "all",
		Stage:       "unused",
	}
}
