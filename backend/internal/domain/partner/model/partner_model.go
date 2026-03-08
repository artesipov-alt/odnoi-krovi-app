package model

import "time"

const (
	RoleClinic  = "CLINIC"
	RoleAdmin   = "ADMIN"
	RoleService = "SERVICE"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusExpired  = "expired"
)

type Partner struct {
	ID          string
	Name        string
	APIKey      string
	Role        string
	Status      string
	Description *string
	LastUsedAt  *time.Time
}
