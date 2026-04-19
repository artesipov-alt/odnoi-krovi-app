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
	PromoCode    string
	ExpiresAt    time.Time
	PlatformName string
	PlatformURL  *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
