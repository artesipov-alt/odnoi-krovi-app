package models

import (
	"time"

	"gorm.io/gorm"
)

// Pet represents a pet in the system
type Pet struct {
	ID                  int             `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	OwnerID             int             `json:"ownerId,omitempty" example:"1"`
	Name                string          `gorm:"size:100;not null" json:"name" example:"Бобик"`
	HasChip             bool            `json:"hasChip" example:"false"`
	ChipNumber          string          `gorm:"size:50" json:"chipNumber,omitempty" example:"123456789"`
	PhotoURL            string          `gorm:"size:255" json:"photoUrl,omitempty" example:"https://example.com/photo.jpg"`
	KnowsBloodGroup     bool            `json:"knowsBloodGroup" example:"false"`
	IsGuideDog          bool            `json:"isGuideDog" example:"false"`
	IsTherapist         bool            `json:"isTherapist" example:"false"`
	Breed               string          `gorm:"size:100" json:"breed,omitempty" example:"Лабрадор"`
	WeightKg            float64         `gorm:"type:numeric" json:"weightKg,omitempty" example:"25.5"`
	AgeYears            int             `json:"ageYears,omitempty" example:"3"`
	AgeMonths           int             `json:"ageMonths,omitempty" example:"6"`
	Sterilized          bool            `json:"sterilized" example:"false"`
	VaccinationDate     *time.Time      `json:"vaccinationDate,omitempty" example:"2023-01-01T12:00:00Z"`
	DewormingDate       *time.Time      `json:"dewormingDate,omitempty" example:"2023-01-01T12:00:00Z"`
	EctoparasiteDate    *time.Time      `json:"ectoparasiteDate,omitempty" example:"2023-01-01T12:00:00Z"`
	LastTransfusionDate *time.Time      `json:"lastTransfusionDate,omitempty" example:"2023-01-01T12:00:00Z"`
	Latitude            float64         `gorm:"type:numeric" json:"latitude,omitempty" example:"55.7558"`
	Longitude           float64         `gorm:"type:numeric" json:"longitude,omitempty" example:"37.6173"`
	LivingCondition     LivingCondition `json:"livingCondition,omitempty" example:"apartment"`
	Gender              Gender          `json:"gender,omitempty" example:"male"`
	Type                PetType         `json:"type,omitempty" example:"dog"`
	BloodGroup          string          `json:"bloodGroup,omitempty" example:"DEA 1.1"`
	DeletedAt           *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}
