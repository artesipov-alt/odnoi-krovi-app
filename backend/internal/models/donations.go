package models

import "time"

// Donation представляет донорство крови в системе
type Donation struct {
	ID         int             `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	DonorPetID *int            `gorm:"column:donor_pet_id" json:"donor_pet_id,omitempty" example:"1"`
	ClinicID   *int            `gorm:"column:clinic_id" json:"clinic_id,omitempty" example:"1"`
	Date       *time.Time      `gorm:"type:date" json:"date,omitempty" example:"2024-01-15"`
	Status     *DonationStatus `gorm:"type:varchar(50)" json:"status,omitempty" example:"completed"`
}
