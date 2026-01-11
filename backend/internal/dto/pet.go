package dto

import (
	"time"
)

// PetHealthDTO represents the health information for a pet
type PetHealthDTO struct {
	ReproductiveStatus    *string    `json:"reproductiveStatus"`
	HealthStatus          *string    `json:"healthStatus"`
	LastDonation          *time.Time `json:"lastDonation"`
	Transfused            *bool      `json:"transfused"`
	Medications           *string    `json:"medications"`
	SurgicalInterventions *string    `json:"surgicalInterventions"`
}

// PetTreatmentDTO represents the treatment information for a pet
type PetTreatmentDTO struct {
	RabiesVaccinationDate     *time.Time `json:"rabiesVaccinationDate"`
	InfectionVaccinationDate  *time.Time `json:"infectionVaccinationDate"`
	EctoparasiteTreatmentDate *time.Time `json:"ectoparasiteTreatmentDate"`
	DewormingDate             *time.Time `json:"dewormingDate"`
}

// PetAnalysisDTO represents the analysis information for a pet
type PetAnalysisDTO struct {
	LeukemiaDate         *time.Time `json:"leukemiaDate"`
	LeukemiaType         *string    `json:"leukemiaType"`
	ImmunodeficiencyDate *time.Time `json:"immunodeficiencyDate"`
	ImmunodeficiencyType *string    `json:"immunodeficiencyType"`
	HemoplasmosisDate    *time.Time `json:"hemoplasmosisDate"`
	HemoplasmosisType    *string    `json:"hemoplasmosisType"`
	BartonellosisDate    *time.Time `json:"bartonellosisDate"`
	BartonellosisType    *string    `json:"bartonellosisType"`
	BabesiosisDate       *time.Time `json:"babesiosisDate"`
	BabesiosisType       *string    `json:"babesiosisType"`
	DirofilariaDate      *time.Time `json:"dirofilariaDate"`
	DirofilariaType      *string    `json:"dirofilariaType"`
	EhrlichiosisDate     *time.Time `json:"ehrlichiosisDate"`
	EhrlichiosisType     *string    `json:"ehrlichiosisType"`
	AnaplasmosisDate     *time.Time `json:"anaplasmosisDate"`
	AnaplasmosisType     *string    `json:"anaplasmosisType"`
}

// PetBonusDTO represents the bonus information for a pet
type PetBonusDTO struct {
	IsArtist      bool `json:"isArtist"`
	IsTherapist   bool `json:"isTherapist"`
	IsFormerDonor bool `json:"isFormerDonor"`
	IsGuideDog    bool `json:"isGuideDog"`
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
	Health          *PetHealthDTO     `json:"health"`
	Treatments      *PetTreatmentDTO  `json:"treatments"`
	Analyses        []*PetAnalysisDTO `json:"analyses"`
	Bonuses         *PetBonusDTO      `json:"bonuses"`
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
	Health          *PetHealthDTO     `json:"health"`
	Treatments      *PetTreatmentDTO  `json:"treatments"`
	Analyses        []*PetAnalysisDTO `json:"analyses"`
	Bonuses         *PetBonusDTO      `json:"bonuses"`
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
