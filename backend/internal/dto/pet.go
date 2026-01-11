package dto

import (
	"time"
)

// PetHealthDTO represents the health information for a pet
type PetHealthDTO struct {
	ReproductiveStatus    *string    `json:"reproductiveStatus,omitempty"`
	HealthStatus          *string    `json:"healthStatus,omitempty"`
	LastDonation          *time.Time `json:"lastDonation,omitempty"`
	Transfused            *bool      `json:"transfused,omitempty"`
	Medications           *string    `json:"medications,omitempty"`
	SurgicalInterventions *string    `json:"surgicalInterventions,omitempty"`
}

// PetTreatmentDTO represents the treatment information for a pet
type PetTreatmentDTO struct {
	RabiesVaccinationDate     *time.Time `json:"rabiesVaccinationDate,omitempty"`
	InfectionVaccinationDate  *time.Time `json:"infectionVaccinationDate,omitempty"`
	EctoparasiteTreatmentDate *time.Time `json:"ectoparasiteTreatmentDate,omitempty"`
	DewormingDate             *time.Time `json:"dewormingDate,omitempty"`
}

// PetAnalysisDTO represents the analysis information for a pet
type PetAnalysisDTO struct {
	LeukemiaDate         *time.Time `json:"leukemiaDate,omitempty"`
	LeukemiaType         *string    `json:"leukemiaType,omitempty"`
	ImmunodeficiencyDate *time.Time `json:"immunodeficiencyDate,omitempty"`
	ImmunodeficiencyType *string    `json:"immunodeficiencyType,omitempty"`
	HemoplasmosisDate    *time.Time `json:"hemoplasmosisDate,omitempty"`
	HemoplasmosisType    *string    `json:"hemoplasmosisType,omitempty"`
	BartonellosisDate    *time.Time `json:"bartonellosisDate,omitempty"`
	BartonellosisType    *string    `json:"bartonellosisType,omitempty"`
	BabesiosisDate       *time.Time `json:"babesiosisDate,omitempty"`
	BabesiosisType       *string    `json:"babesiosisType,omitempty"`
	DirofilariaDate      *time.Time `json:"dirofilariaDate,omitempty"`
	DirofilariaType      *string    `json:"dirofilariaType,omitempty"`
	EhrlichiosisDate     *time.Time `json:"ehrlichiosisDate,omitempty"`
	EhrlichiosisType     *string    `json:"ehrlichiosisType,omitempty"`
	AnaplasmosisDate     *time.Time `json:"anaplasmosisDate,omitempty"`
	AnaplasmosisType     *string    `json:"anaplasmosisType,omitempty"`
}

// PetBonusDTO represents the bonus information for a pet
type PetBonusDTO struct {
	IsArtist      bool `json:"isArtist,omitempty"`
	IsTherapist   bool `json:"isTherapist,omitempty"`
	IsFormerDonor bool `json:"isFormerDonor,omitempty"`
	IsGuideDog    bool `json:"isGuideDog,omitempty"`
}

// PetCreate represents the structure for creating a new pet
type PetCreate struct {
	Name            string            `json:"name" validate:"required,min=1,max=100"`
	ChipNumber      string            `json:"chipNumber,omitempty" validate:"omitempty,len=15"`
	PhotoURL        string            `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	BreedID         int               `json:"breedId,omitempty" validate:"omitempty,min=1"`
	WeightKg        float64           `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        int               `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       int               `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	BirthDate       *time.Time        `json:"birthDate,omitempty"`
	LivingCondition string            `json:"livingCondition,omitempty"`
	Gender          string            `json:"gender,omitempty"`
	Type            string            `json:"type" validate:"required"`
	BloodGroup      string            `json:"bloodGroup,omitempty"`
	PetStatus       string            `json:"petStatus" validate:"required"`
	Health          *PetHealthDTO     `json:"health,omitempty"`
	Treatments      *PetTreatmentDTO  `json:"treatments,omitempty"`
	Analyses        []*PetAnalysisDTO `json:"analyses,omitempty"`
	Bonuses         *PetBonusDTO      `json:"bonuses,omitempty"`
}

// PetUpdate represents the structure for updating an existing pet
type PetUpdate struct {
	Name            *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	ChipNumber      *string           `json:"chipNumber,omitempty" validate:"omitempty,len=15"`
	PhotoURL        *string           `json:"photoUrl,omitempty" validate:"omitempty,url,max=255"`
	BreedID         *int              `json:"breedId,omitempty" validate:"omitempty,min=1"`
	WeightKg        *float64          `json:"weightKg,omitempty" validate:"omitempty,min=0"`
	AgeYears        *int              `json:"ageYears,omitempty" validate:"omitempty,min=0"`
	AgeMonths       *int              `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11"`
	BirthDate       *time.Time        `json:"birthDate,omitempty"`
	LivingCondition *string           `json:"livingCondition,omitempty"`
	Gender          *string           `json:"gender,omitempty"`
	Type            *string           `json:"type,omitempty"`
	BloodGroup      *string           `json:"bloodGroup,omitempty"`
	PetStatus       *string           `json:"petStatus,omitempty"`
	Health          *PetHealthDTO     `json:"health,omitempty"`
	Treatments      *PetTreatmentDTO  `json:"treatments,omitempty"`
	Analyses        []*PetAnalysisDTO `json:"analyses,omitempty"`
	Bonuses         *PetBonusDTO      `json:"bonuses,omitempty"`
}

type PetResponseDTO struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	ChipNumber      string            `json:"chipNumber,omitempty"`
	PhotoURL        string            `json:"photoUrl,omitempty"`
	BreedID         int               `json:"breedId,omitempty"`
	WeightKg        float64           `json:"weightKg,omitempty"`
	AgeYears        int               `json:"ageYears,omitempty"`
	AgeMonths       int               `json:"ageMonths,omitempty"`
	BirthDate       *time.Time        `json:"birthDate,omitempty"`
	LivingCondition string            `json:"livingCondition,omitempty"`
	Gender          string            `json:"gender,omitempty"`
	Type            string            `json:"type"`
	BloodGroup      string            `json:"bloodGroup,omitempty"`
	PetStatus       string            `json:"petStatus"`
	Health          *PetHealthDTO     `json:"health"`
	Treatments      *PetTreatmentDTO  `json:"treatments"`
	Analyses        []*PetAnalysisDTO `json:"analyses"`
	Bonuses         *PetBonusDTO      `json:"bonuses"`
	CreatedAt       string            `json:"createdAt"`
	UpdatedAt       string            `json:"updatedAt"`
}
