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
	Amount          int32
	DonorBloodGroup string
	PhotoURLs       []string
}

func (r *Recipient) AddMatchingDonor(pet *petmodel.Pet) {
	r.MatchingDonors = append(r.MatchingDonors, MatchingDonorReadModel{
		PetName:         pet.Name,
		DonorBloodGroup: *pet.BloodGroupName,
		PhotoURLs:       pet.PhotoURLs,
		Amount:          CalculateDonationAmount(string(pet.Type), pet.WeightKg),
	})
}

func (r *Recipient) SetDefaultPrefs(compensationType usermodel.CompensationType, taxiCompensation bool) {
	r.DefaultDonorPrefs = &DefaultDonorPrefs{
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Bonuses:          []string{},
	}
}

// calculateDonationAmount вычисляет максимальный объем донации крови для питомца (до 20% циркулирующей крови, но не более лимита)
// Для собак: не более 17.6 мл/кг
// Для кошек: не более 13.2 мл/кг
func CalculateDonationAmount(petType string, weightKg float64) int32 {
	var limitPerKg float64
	switch petType {
	case "dog":
		limitPerKg = 17.6
	case "cat":
		limitPerKg = 13.2
	default:
		return 0
	}
	amount := limitPerKg * weightKg
	return int32(amount)
}
