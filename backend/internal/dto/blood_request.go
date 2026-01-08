package dto

import "time"

// BloodSearchPetRequest represents a request to add a pet to the blood search pool
type BloodSearchPetRequest struct {
	PetID                  string   `json:"petId" validate:"required"`
	BloodVolumeNeeded      int32    `json:"bloodVolumeNeeded" validate:"required,gt=0"`
	BloodVolumeReserved    int32    `json:"bloodVolumeReserved"`
	Regions                []int32  `json:"regions" validate:"required,min=1"`
	SmallPetsNotifyAllowed bool     `json:"smallPetsNotifyAllowed"`
	Description            string   `json:"description"`
	PhotoUrls              []string `json:"photoUrls"`
	BloodGroupIds          []string `json:"bloodGroupIds"`
	BloodComponentIds      []int    `json:"bloodComponentIds"`
}

// BloodSearchPetResponse represents the response after creating a blood search request
type BloodSearchPetResponse struct {
	ID     string `json:"id"`
	PetID  string `json:"petId"`
	Status string `json:"status"`
}

// BloodSearchFilterRequest represents filters for searching blood requests
type BloodSearchFilterRequest struct {
	PetID  string `json:"petId"`
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// BloodSearchRequestDTO represents a DTO for a blood search request
type BloodSearchRequestDTO struct {
	ID                     string    `json:"id"`
	PetID                  string    `json:"petId"`
	BloodVolumeNeeded      int32     `json:"bloodVolumeNeeded"`
	BloodVolumeReserved    int32     `json:"bloodVolumeReserved"`
	Regions                []int32   `json:"regions"`
	SmallPetsNotifyAllowed bool      `json:"smallPetsNotifyAllowed"`
	Description            string    `json:"description"`
	PhotoUrls              []string  `json:"photoUrls"`
	BloodGroupIds          []string  `json:"bloodGroupIds"`
	BloodComponentIds      []int     `json:"bloodComponentIds"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

// BloodSearchPetsResponse represents a list of blood search requests
type BloodSearchPetsResponse struct {
	Requests []BloodSearchRequestDTO `json:"requests"`
}
