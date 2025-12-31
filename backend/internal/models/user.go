package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UserPrefix определяет тип для префиксов пользователя
type UserPrefix string

const (
	// UserIDPrefix — основной префикс для пользователей
	UserIDPrefix UserPrefix = "USR"
)

// Представление пользователя в системе
type User struct {
	//Сигнатура ID пользователя включает в себя префикс пользователя, год и четырёхзначный номер
	ID               string   `gorm:"primaryKey;size:15" json:"id" example:"USR-25-000001"`
	TelegramID       int64    `gorm:"not null" json:"telegramId" example:"123456789"`
	FullName         string   `gorm:"size:255" json:"fullName,omitempty" example:"Иван Иванов"`
	Phone            string   `gorm:"size:20" json:"phone,omitempty" example:"+79991234567"`
	Email            string   `gorm:"size:255" json:"email,omitempty" example:"user@example.com"`
	OrganizationName string   `gorm:"size:255" json:"organizationName,omitempty" example:"ООО Ромашка"`
	ConsentPD        bool     `json:"consentPd" example:"true"`
	OnBoarding       bool     `json:"onBoarding" example:"false"`
	AllowGeo         bool     `json:"allowGeo" example:"true"`
	LocationID       int      `json:"locationId,omitempty" example:"1"`
	Role             UserRole `json:"role,omitempty" example:"user"`
	// Дополнительные данные базы
	CreatedAt time.Time       `json:"createdAt,omitempty" swaggerignore:"true" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time       `json:"updatedAt,omitempty" swaggerignore:"true" example:"2023-01-01T00:00:00Z"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// UserRole представляет роль пользователя в системе
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// Generate создает префикс для сущности пользователя
func (e UserPrefix) Generate(sequenceNum int) string {
	year := time.Now().Year() % 100
	return fmt.Sprintf("%s-%02d-%06d", e, year, sequenceNum)
}

// IsValid проверяет валидность префикса пользователя
func (e UserPrefix) IsValid() bool {
	return e == UserIDPrefix
}

// BeforeCreate хук для генерации ID
func (v *User) BeforeCreate(tx *gorm.DB) error {
	var nextVal int
	if err := tx.Raw("SELECT nextval('user_id_seq')").Scan(&nextVal).Error; err != nil {
		return err
	}
	v.ID = UserIDPrefix.Generate(nextVal)
	return nil
}
