package models

import "time"

// BloodStock представляет запас крови в системе
type BloodStock struct {
	ID             int              `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	ClinicID       *int             `gorm:"index" json:"clinicId,omitempty" example:"1"`
	PetType        PetType          `gorm:"type:varchar(50)" json:"petType" example:"dog"`
	VolumeML       *int             `json:"volumeMl,omitempty" example:"500"`
	PriceRub       *float64         `gorm:"type:numeric" json:"priceRub,omitempty" example:"5000.00"`
	ExpirationDate *time.Time       `gorm:"type:date" json:"expirationDate,omitempty" example:"2024-12-31"`
	Status         BloodStockStatus `gorm:"type:varchar(50);default:'active'" json:"status,omitempty" example:"active"`
	BloodTypeID    int              `gorm:"index" json:"bloodTypeId" example:"1"`
}
