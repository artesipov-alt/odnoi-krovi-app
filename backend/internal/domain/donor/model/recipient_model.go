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
	PhotoURLs            []string
	BloodGroupName       string
	PrioritySearch       bool
	Status               string
	MatchingDonors       []MatchingDonorReadModel
	DefaultDonorPrefs    *DefaultDonorPrefs
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
	PetType         petmodel.PetType
	DonorBloodGroup string
	PhotoURLs       []string
}

func (r *Recipient) AddMatchingDonor(donorName, donorBloodGroup string, donorType petmodel.PetType, photoURLs []string) {
	r.MatchingDonors = append(r.MatchingDonors, MatchingDonorReadModel{
		PetName:         donorName,
		PetType:         donorType,
		DonorBloodGroup: donorBloodGroup,
		PhotoURLs:       photoURLs,
	})
}

func (r *Recipient) SetDefaultPrefs(compensationType usermodel.CompensationType, taxiCompensation bool) {
	r.DefaultDonorPrefs = &DefaultDonorPrefs{
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Bonuses:          []string{},
	}
}
