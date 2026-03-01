package model

import (
	"errors"
	"time"

	pet "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// UserRole represents user role types
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User представляет доменную модель пользователя
type User struct {
	ID               string
	TelegramID       int64
	FullName         string
	Phone            string
	Email            string
	PhotoURLs        []string
	OrganizationName string
	ConsentPd        bool
	OnBoarding       []string
	AllowGeo         bool
	LocationID       *string
	Role             string
	Pets             []*pet.Pet
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// NewUser creates a new User aggregate with validation
func NewUser(
	telegramID int64,
	fullName string,
	phone string,
	email string,
	role UserRole,
	consentPd bool,
	locationID *string,
) (*User, error) {
	// Validation
	if fullName == "" {
		return nil, errors.New("full name is required")
	}
	if len(fullName) > 100 {
		return nil, errors.New("full name must be less than 100 characters")
	}
	if role == "" {
		role = RoleUser // default role
	}
	if role != RoleUser && role != RoleAdmin {
		return nil, errors.New("invalid user role")
	}
	if email != "" {
		// Basic email validation could be added here
		if len(email) > 100 {
			return nil, errors.New("email must be less than 100 characters")
		}
	}
	if phone != "" && len(phone) > 20 {
		return nil, errors.New("phone must be less than 20 characters")
	}

	user := &User{
		TelegramID: telegramID,
		FullName:   fullName,
		Phone:      phone,
		Email:      email,
		Role:       string(role),
		ConsentPd:  consentPd,
		LocationID: locationID,
		PhotoURLs:  []string{},
		OnBoarding: []string{},
		Pets:       []*pet.Pet{},
	}

	return user, nil
}
