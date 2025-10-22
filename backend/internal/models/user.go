package models

import (
	"time"
)

// User представляет пользователя в системе
type User struct {
	ID               int        `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	TelegramID       int64      `gorm:"not null" json:"telegram_id" example:"123456789"`
	FullName         string     `gorm:"size:255" json:"full_name,omitempty" example:"Иван Иванов"`
	Phone            string     `gorm:"size:20" json:"phone,omitempty" example:"+79991234567"`
	Email            string     `gorm:"size:255" json:"email,omitempty" example:"user@example.com"`
	OrganizationName string     `gorm:"size:255" json:"organization_name,omitempty" example:"ООО Ромашка"`
	CreatedAt        time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	ConsentPD        bool       `json:"consent_pd" example:"true"`
	LocationID       int        `json:"location_id,omitempty" example:"1"`
	Role             string     `json:"role,omitempty" example:"user"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty" example:"2023-01-02T12:00:00Z"`
}
