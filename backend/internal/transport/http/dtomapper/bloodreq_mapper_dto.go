// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
)

// BloodRequestMapper handles conversions between domain BloodRequest model and DTOs.
type BloodRequestMapper struct {
	storage filestorage.Repository
}

// NewBloodRequestMapper creates a new BloodRequestMapper instance.
func NewBloodRequestMapper(storage filestorage.Repository) *BloodRequestMapper {
	return &BloodRequestMapper{
		storage: storage,
	}
}

// ToResponse converts a domain BloodRequest model to a DTO.
func (m *BloodRequestMapper) ToResponse(req *model.BloodRequest, suitableDonors *int) dto.BloodRequestDetail {
	if req == nil {
		return dto.BloodRequestDetail{}
	}

	suitableDonorsCount := 0
	if suitableDonors != nil {
		suitableDonorsCount = *suitableDonors
	}

	var donorApplications []dto.DonorApplication
	if len(req.DonorApplications) != 0 {
		donorApplications = make([]dto.DonorApplication, 0, len(req.DonorApplications))
		for _, app := range req.DonorApplications {
			// Convert WarnFactors from []string to []RestrictionFactor with descriptions
			var warnFactors []dto.RestrictionFactor
			for _, code := range app.WarnFactors {
				desc := petmodel.GetFactorDescription(petmodel.FactorCode(code))
				factor := dto.RestrictionFactor{
					Code:           code,
					Description:    desc.Description,
					SubDescription: desc.SubDescription,
				}
				warnFactors = append(warnFactors, factor)
			}

			donorApplications = append(donorApplications, dto.DonorApplication{
				ID:               app.ID,
				RequestID:        app.RequestID,
				DonorID:          app.DonorID,
				DonorName:        app.DonorName,
				DonorPhotos:      m.storage.BuildPhotoURLs(app.DonorPhotos, *req.UpdatedAt),
				DonorBloodGroup:  app.DonorBloodGroup,
				Amount:           app.Amount,
				WarnFactors:      warnFactors,
				CompensationType: app.CompensationType,
				TaxiCompensation: app.TaxiCompensation,
				Status:           string(app.Status),
				CreatedAt:        app.CreatedAt,
				UpdatedAt:        app.UpdatedAt,
			})
		}
	}

	return dto.BloodRequestDetail{
		ID:                       req.ID,
		PetID:                    req.PetID,
		BloodVolumeNeeded:        req.BloodVolumeNeeded,
		BloodVolumeReserved:      req.BloodVolumeReserved,
		Regions:                  req.Regions,
		SmallPetsNotifyAllowed:   req.SmallPetsNotifyAllowed,
		Description:              req.Description,
		PhotoURLs:                m.storage.BuildPhotoURLs(req.PhotoURLs, *req.UpdatedAt),
		BloodGroupNames:          req.BloodGroupNames,
		BloodComponentIDs:        req.BloodComponentIDs,
		OnBoarding:               req.OnBoarding,
		Status:                   string(req.Status),
		SuitableDonors:           suitableDonorsCount,
		PrioritySearch:           req.PrioritySearch,
		IncludeUnknownBloodGroup: req.IncludeUnknownBloodGroup,
		Responses:                donorApplications,
		CreatedAt:                req.CreatedAt,
		UpdatedAt:                req.UpdatedAt,
		DeletedAt:                req.DeletedAt,
	}
}

// ToResponseSlice converts a slice of domain BloodRequest models to DTOs.
func (m *BloodRequestMapper) ToResponseSlice(reqs []*model.BloodRequest) []dto.BloodRequestDetail {
	if reqs == nil {
		return nil
	}
	dtos := make([]dto.BloodRequestDetail, len(reqs))
	for i, req := range reqs {
		dtos[i] = m.ToResponse(req, nil)
	}
	return dtos
}

// FromCreate converts a CreateBloodRequestBody DTO to a domain BloodRequest model using the constructor.
func (m *BloodRequestMapper) FromCreate(body dto.CreateBloodRequestBody) *model.BloodRequest {
	req := model.NewBloodRequest(
		body.PetID,
		body.BloodVolumeNeeded,
		body.Regions,
	)

	// Set additional fields from DTO
	req.SmallPetsNotifyAllowed = body.SmallPetsNotifyAllowed
	req.Description = body.Description
	req.BloodGroupNames = body.BloodGroupNames
	req.BloodComponentIDs = body.BloodComponentIDs
	req.PrioritySearch = body.PrioritySearch
	req.IncludeUnknownBloodGroup = body.IncludeUnknownBloodGroup

	return req
}
