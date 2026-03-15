package model

import petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"

// Recipient представляет модель чтения для списка реципиентов
type Recipient struct {
	ID                   string
	PetID                string
	PetName              string
	PetType              petmodel.PetType
	BloodVolumeRemaining int32
	PhotoURLs            []string
	BloodGroupName       string
	PrioritySearch       bool
	Status               string
	MatchingDonors       []MatchingDonorReadModel
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
	DonorBloodGroup string
	PhotoURLs       []string
}
