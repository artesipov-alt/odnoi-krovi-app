package models

import (
	"time"

	"gorm.io/gorm"
)

// User представляет пользователя в системе
type User struct {
	ID               int             `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	TelegramID       int64           `gorm:"not null" json:"telegramId" example:"123456789"`
	FullName         string          `gorm:"size:255" json:"fullName,omitempty" example:"Иван Иванов"`
	Phone            string          `gorm:"size:20" json:"phone,omitempty" example:"+79991234567"`
	Email            string          `gorm:"size:255" json:"email,omitempty" example:"user@example.com"`
	OrganizationName string          `gorm:"size:255" json:"organizationName,omitempty" example:"ООО Ромашка"`
	CreatedAt        time.Time       `json:"createdAt" example:"2023-01-01T12:00:00Z"`
	ConsentPD        bool            `json:"consentPd" example:"true"`
	OnBoarding       bool            `json:"onBoarding" example:"false"`
	AllowGeo         bool            `json:"allowGeo" example:"true"`
	LocationID       int             `json:"locationId,omitempty" example:"1"`
	Role             UserRole        `json:"role,omitempty" example:"user"`
	DeletedAt        *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}
