package model

import (
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

// Recipient представляет модель чтения реципиент
type Recipient struct {
	ID                   string
	PetID                string
	PetName              string
	PetType              petmodel.PetType
	OwnerName            string
	SearchRegions        []string
	SearchingBloodNames  []string
	BloodVolumeRemaining int32
	BloodVolumeNeeded    int32
	BloodVolumeReserved  int32
	PhotoURLs            []string
	BloodGroupName       string
	PrioritySearch       bool
	Status               string
	MatchingDonors       []MatchingDonorReadModel
	DefaultDonorPrefs    *DefaultDonorPrefs
	AdvancedInfo         *AdvancedInfo
}

type AdvancedInfo struct {
	Description string
	PhotoURLs   []string
}

type DefaultDonorPrefs struct {
	CompensationType usermodel.CompensationType
	Bonuses          []string
	TaxiCompensation bool
}

// DonorPreloadFilter представляет параметры для предзагрузки связанных данных
type DonorPreloadFilter struct {
	Status string
	Limit  int
	Offset int
}

// MatchingDonorReadModel представляет модель чтения для подходящего донора
type MatchingDonorReadModel struct {
	PetID           string
	PetName         string
	Amount          int32
	DonorBloodGroup string
	PhotoURLs       []string
}

func (r *Recipient) AddMatchingDonor(pet *petmodel.Pet) {
	r.MatchingDonors = append(r.MatchingDonors, MatchingDonorReadModel{
		PetName:         pet.Name,
		PetID:           pet.ID,
		DonorBloodGroup: *pet.BloodGroupName,
		PhotoURLs:       pet.PhotoURLs,
		Amount:          pet.CalculateDonationAmount(),
	})
}

func (r *Recipient) SetDefaultPrefs(compensationType usermodel.CompensationType, taxiCompensation bool) {
	r.DefaultDonorPrefs = &DefaultDonorPrefs{
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Bonuses:          []string{},
	}
}
