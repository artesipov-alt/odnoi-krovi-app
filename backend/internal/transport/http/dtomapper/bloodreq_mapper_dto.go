// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
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

func (m *BloodRequestMapper) ToResponse(req *model.BloodRequestWithApplications, suitableDonors *int) dto.BloodRequestDetail {
	if req == nil {
		return dto.BloodRequestDetail{}
	}

	suitableDonorsCount := 0
	if suitableDonors != nil {
		suitableDonorsCount = *suitableDonors
	}

	var donorApplications []dto.DonorApplication
	var acceptedDonorApplications []dto.DonorApplication
	var completedDonations []dto.DonorApplication
	if len(req.DonorApplications) != 0 {
		donorApplications = make([]dto.DonorApplication, 0, len(req.DonorApplications))
		acceptedDonorApplications = make([]dto.DonorApplication, 0, len(req.DonorApplications))
		completedDonations = make([]dto.DonorApplication, 0, len(req.DonorApplications))
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

			dtoApp := dto.DonorApplication{
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
				IsConfirmed:      app.IsConfirmed,
				CreatedAt:        app.CreatedAt,
				UpdatedAt:        app.UpdatedAt,
			}

			if app.Status == donormodel.DonorResponseStatusPending {
				donorApplications = append(donorApplications, dtoApp)
			} else if app.Status == donormodel.DonorResponseStatusCompleted && app.IsConfirmed {
				completedDonations = append(completedDonations, dtoApp)
			} else if app.IsConfirmed != true && app.Status != donormodel.DonorResponseStatusPending {
				acceptedDonorApplications = append(acceptedDonorApplications, dtoApp)
			}

		}
	}

	return dto.BloodRequestDetail{
		ID:                       req.ID,
		PetID:                    req.PetID,
		BloodVolumeNeeded:        req.BloodVolumeNeeded,
		BloodVolumeReserved:      req.BloodVolumeReserved,
		BloodVolumeDonated:       req.BloodVolumeDonated,
		Regions:                  req.Regions,
		SmallPetsNotifyAllowed:   req.SmallPetsNotifyAllowed,
		Description:              req.AdvancedInfo.Description,
		PhotoURLs:                m.storage.BuildPhotoURLs(req.AdvancedInfo.PhotoURLs, *req.UpdatedAt),
		BloodGroupNames:          req.BloodGroupNames,
		BloodComponentIDs:        req.BloodComponentIDs,
		OnBoarding:               req.OnBoarding,
		Status:                   string(req.Status),
		PrioritySearch:           req.PrioritySearch,
		IncludeUnknownBloodGroup: req.IncludeUnknownBloodGroup,
		Responses:                donorApplications,
		AcceptedDonors:           acceptedDonorApplications,
		CompletedDonations:       completedDonations,
		SuitableDonors:           suitableDonorsCount,
		CreatedAt:                req.CreatedAt,
		UpdatedAt:                req.UpdatedAt,
		DeletedAt:                req.DeletedAt,
	}
}

// ToResponseSlice converts a slice of domain BloodRequest models to DTOs.
func (m *BloodRequestMapper) ToResponseSlice(reqs []*model.BloodRequestWithApplications) []dto.BloodRequestDetail {
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
	req.AdvancedInfo.Description = body.Description
	req.BloodGroupNames = body.BloodGroupNames
	req.BloodComponentIDs = body.BloodComponentIDs
	req.PrioritySearch = body.PrioritySearch
	req.IncludeUnknownBloodGroup = body.IncludeUnknownBloodGroup

	return req
}
