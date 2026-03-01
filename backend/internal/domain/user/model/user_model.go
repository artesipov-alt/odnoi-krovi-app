package model

import (
	"time"

	pet "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
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
