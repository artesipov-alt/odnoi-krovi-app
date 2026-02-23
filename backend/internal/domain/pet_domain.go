package domain

import (
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
	WeightKg           float64
	Gender             Gender
	BirthDate          *time.Time
	ChipNumber         string
	PhotoURLs          []string
	LivingCondition    LivingCondition
	ReproductiveStatus ReproductiveStatus
	DonorRestrictions  []string
	OwnerID            string
	BreedRefID         *string
	BloodGroupRefID    *string
	Health             *PetHealth
	Treatments         *PetTreatment
	Analyses           []*PetAnalysis
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
	DeletedAt          *time.Time
}

// PetHealth представляет здоровье питомца
type PetHealth struct {
	HealthStatus          HealthStatus
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
