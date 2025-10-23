package models

import "time"

// BloodStock представляет запас крови в системе
type BloodStock struct {
	ID             int              `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	ClinicID       *int             `gorm:"index" json:"clinic_id,omitempty" example:"1"`
	PetType        PetType          `gorm:"type:varchar(50)" json:"pet_type" example:"dog"`
	VolumeML       *int             `json:"volume_ml,omitempty" example:"500"`
	PriceRub       *float64         `gorm:"type:numeric" json:"price_rub,omitempty" example:"5000.00"`
	ExpirationDate *time.Time       `gorm:"type:date" json:"expiration_date,omitempty" example:"2024-12-31"`
	Status         BloodStockStatus `gorm:"type:varchar(50);default:'active'" json:"status,omitempty" example:"active"`
	BloodTypeID    int              `gorm:"index" json:"blood_type_id" example:"1"`
}
