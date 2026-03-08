package model

import "time"

type Partner struct {
	ID          string
	Name        string
	APIKey      string
	Role        string
	Status      string
	Description *string
	LastUsedAt  *time.Time
}
