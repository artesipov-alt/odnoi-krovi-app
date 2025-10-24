package models

import (
	"time"

	"gorm.io/gorm"
)

// Pet represents a pet in the system
type Pet struct {
	ID                  int             `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	OwnerID             int             `json:"owner_id,omitempty" example:"1"`
	Name                string          `gorm:"size:100;not null" json:"name" example:"Бобик"`
	HasChip             bool            `json:"has_chip" example:"false"`
	ChipNumber          string          `gorm:"size:50" json:"chip_number,omitempty" example:"123456789"`
	PhotoURL            string          `gorm:"size:255" json:"photo_url,omitempty" example:"https://example.com/photo.jpg"`
	KnowsBloodGroup     bool            `json:"knows_blood_group" example:"false"`
	IsGuideDog          bool            `json:"is_guide_dog" example:"false"`
	IsTherapist         bool            `json:"is_therapist" example:"false"`
	Breed               string          `gorm:"size:100" json:"breed,omitempty" example:"Лабрадор"`
	WeightKg            float64         `gorm:"type:numeric" json:"weight_kg,omitempty" example:"25.5"`
	AgeYears            int             `json:"age_years,omitempty" example:"3"`
	AgeMonths           int             `json:"age_months,omitempty" example:"6"`
	Sterilized          bool            `json:"sterilized" example:"false"`
	VaccinationDate     *time.Time      `json:"vaccination_date,omitempty" example:"2023-01-01T12:00:00Z"`
	DewormingDate       *time.Time      `json:"deworming_date,omitempty" example:"2023-01-01T12:00:00Z"`
	EctoparasiteDate    *time.Time      `json:"ectoparasite_date,omitempty" example:"2023-01-01T12:00:00Z"`
	LastTransfusionDate *time.Time      `json:"last_transfusion_date,omitempty" example:"2023-01-01T12:00:00Z"`
	Latitude            float64         `gorm:"type:numeric" json:"latitude,omitempty" example:"55.7558"`
	Longitude           float64         `gorm:"type:numeric" json:"longitude,omitempty" example:"37.6173"`
	LivingCondition     LivingCondition `json:"living_condition,omitempty" example:"apartment"`
	Gender              Gender          `json:"gender,omitempty" example:"male"`
	Type                PetType         `json:"type,omitempty" example:"dog"`
	BloodGroup          string          `json:"blood_group,omitempty" example:"DEA 1.1"`
	DeletedAt           *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}
