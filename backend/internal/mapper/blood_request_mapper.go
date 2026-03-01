// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
)

// BloodRequestMapper handles conversions between domain BloodRequest model and DTOs.
type BloodRequestMapper struct{}

// NewBloodRequestMapper creates a new BloodRequestMapper instance.
func NewBloodRequestMapper() *BloodRequestMapper {
	return &BloodRequestMapper{}
}

// ToResponse converts a domain BloodRequest model to a DTO.
func (m *BloodRequestMapper) ToResponse(req *model.BloodRequest, suitableDonors *int) dto.BloodSearchRequest {
	if req == nil {
		return dto.BloodSearchRequest{}
	}

	var createdAt, updatedAt, deletedAt *time.Time

	if !req.CreatedAt.IsZero() {
		createdAt = &req.CreatedAt
	}
	if !req.UpdatedAt.IsZero() {
		updatedAt = &req.UpdatedAt
	}
	deletedAt = req.DeletedAt

	suitableDonorsCount := 0
	if suitableDonors != nil {
		suitableDonorsCount = *suitableDonors
	}

	return dto.BloodSearchRequest{
		ID:                     req.ID,
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoURLs,
		BloodGroupNames:        req.BloodGroupNames,
		BloodComponentIds:      req.BloodComponentIDs,
		OnBoarding:             req.OnBoarding,
		Status:                 dto.BloodSearchRequestStatus(req.Status),
		SuitableDonors:         suitableDonorsCount,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		DeletedAt:              deletedAt,
	}
}

// ToResponseSlice converts a slice of domain BloodRequest models to DTOs.
func (m *BloodRequestMapper) ToResponseSlice(reqs []*model.BloodRequest) []dto.BloodSearchRequest {
	if reqs == nil {
		return nil
	}
	dtos := make([]dto.BloodSearchRequest, len(reqs))
	for i, req := range reqs {
		dtos[i] = m.ToResponse(req, nil)
	}
	return dtos
}

// FromCreate converts a CreateBloodSearchRequest DTO to a domain BloodRequest model.
func (m *BloodRequestMapper) FromCreate(dto dto.CreateBloodSearchRequest) *model.BloodRequest {
	return &model.BloodRequest{
		PetID:                  dto.PetID,
		BloodVolumeNeeded:      dto.BloodVolumeNeeded,
		BloodVolumeReserved:    0,
		Regions:                dto.Regions,
		SmallPetsNotifyAllowed: dto.SmallPetsNotifyAllowed,
		Description:            dto.Description,
		PhotoURLs:              nil,
		BloodGroupNames:        dto.BloodGroupNames,
		BloodComponentIDs:      dto.BloodComponentIds,
		OnBoarding:             nil,
		Status:                 model.BloodRequestStatusActive,
	}
}
