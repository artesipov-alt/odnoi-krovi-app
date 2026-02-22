package domain

import (
	"errors"
	"time"
)

// PetStatus представляет статус питомца
type PetStatus string

const (
	PetStatusNone       PetStatus = "none"
	PetStatusDonor      PetStatus = "donor"
	PetStatusRecipient  PetStatus = "recipient"
	PetStatusBloodFound PetStatus = "blood_found"
)

// PetType представляет тип животного
type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

// Gender представляет пол питомца
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// LivingCondition представляет условия проживания
type LivingCondition string

const (
	LivingConditionIndoor       LivingCondition = "indoor"
	LivingConditionLeashWalking LivingCondition = "leash_walking"
	LivingConditionSelfOutdoor  LivingCondition = "self_outdoor"
)

// ReproductiveStatus представляет репродуктивный статус
type ReproductiveStatus string

const (
	ReproductiveStatusPregnancy ReproductiveStatus = "pregnancy"
	ReproductiveStatusLactation ReproductiveStatus = "lactation"
	ReproductiveStatusEstrus    ReproductiveStatus = "estrus"
)

// HealthStatus представляет состояние здоровья
type HealthStatus string

const (
	HealthStatusHealthy HealthStatus = "healthy"
	HealthStatusIll     HealthStatus = "ill"
	HealthStatusUnknown HealthStatus = "unknown"
)

// Pet представляет доменную модель питомца
type Pet struct {
	ID                 string
	Name               string
	Type               PetType
	WeightKg           *float64
	Gender             *Gender
	BirthDate          *time.Time
	ChipNumber         *string
	PhotoURLs          []string
	BreedID            *string
	UserID             *string
	LivingCondition    *LivingCondition
	ReproductiveStatus *ReproductiveStatus
	DonorRestrictions  []string
	BloodGroupID       *string
	PetStatus          PetStatus
	Health             *PetHealth
	Treatments         *PetTreatment
	Analyses           []*PetAnalysis
	Bonuses            *PetBonus
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}

// PetHealth представляет здоровье питомца
type PetHealth struct {
	HealthStatus          *HealthStatus
	LastDonation          *time.Time
	Transfused            *bool
	Medications           *string
	SurgicalInterventions *string
}

// PetTreatment представляет лечение питомца
type PetTreatment struct {
	RabiesVaccinationDate     *time.Time
	InfectionVaccinationDate  *time.Time
	EctoparasiteTreatmentDate *time.Time
	DewormingDate             *time.Time
}

// PetAnalysis представляет анализ питомца
type PetAnalysis struct {
	ID           string
	AnalysisName string
	AnalysisType string
	AnalysisDate *time.Time
}

// PetBonus представляет бонусы питомца
type PetBonus struct {
	IsArtist      bool
	IsTherapist   bool
	IsFormerDonor bool
	IsGuideDog    bool
}

// IsEligibleForDonation проверяет, подходит ли питомец для донорства
func (p *Pet) IsEligibleForDonation() error {
	if p.PetStatus != PetStatusDonor {
		return errors.New("pet is not a donor")
	}
	if p.Health != nil && p.Health.HealthStatus != nil && *p.Health.HealthStatus == HealthStatusIll {
		return errors.New("pet is ill")
	}
	if p.ReproductiveStatus != nil && (*p.ReproductiveStatus == ReproductiveStatusPregnancy || *p.ReproductiveStatus == ReproductiveStatusLactation) {
		return errors.New("pet is in reproductive status")
	}
	// Дополнительные проверки по весу, возрасту и т.д.
	if p.WeightKg != nil && *p.WeightKg < 5 {
		return errors.New("pet weight is too low")
	}
	return nil
}

// CalculateAge рассчитывает возраст питомца в годах и месяцах
func (p *Pet) CalculateAge() (years int, months int) {
	if p.BirthDate == nil {
		return 0, 0
	}
	now := time.Now()
	years = now.Year() - p.BirthDate.Year()
	if now.YearDay() < p.BirthDate.YearDay() {
		years--
	}
	months = int(now.Sub(*p.BirthDate).Hours() / 24 / 30)
	return years, months % 12
}
