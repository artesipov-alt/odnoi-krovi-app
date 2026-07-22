// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"math"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/query"
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

// DonorResponseToApplication конвертирует PotentialDonor в DonorApplication DTO
// для отображения в списке потенциальных доноров.
func (m *BloodRequestMapper) DonorResponseToApplication(pd *model.PotentialDonor) dto.DonorApplication {
	if pd == nil || pd.Pet == nil {
		return dto.DonorApplication{}
	}

	pet := pd.Pet

	var warnFactors []dto.RestrictionFactor
	for _, code := range pet.WarnFactors {
		desc := petmodel.GetFactorDescription(petmodel.FactorCode(code))
		warnFactors = append(warnFactors, dto.RestrictionFactor{
			Code:           code,
			Description:    desc.Description,
			SubDescription: desc.SubDescription,
		})
	}

	return dto.DonorApplication{
		DonorID:          pet.ID,
		DonorName:        pet.Name,
		DonorBloodGroup:  pet.BloodGroupName,
		DonorPhotos:      m.storage.BuildPhotoURLs(pet.PhotoURLs, *pet.UpdatedAt),
		Amount:           pet.CalculateDonationAmount(),
		WarnFactors:      warnFactors,
		CompensationType: string(pd.CompensationType),
		TaxiCompensation: pd.TaxiCompensation,
		Status:           "",
		IsConfirmed:      false,
	}
}

// DonorResponseToSelected конвертирует DonorResponse (созданный при выборе донора)
// в DonorApplication DTO для ответа на запрос select-donor.
func (m *BloodRequestMapper) DonorResponseToSelected(resp *donormodel.DonorResponse) dto.DonorApplication {
	if resp == nil {
		return dto.DonorApplication{}
	}

	var warnFactors []dto.RestrictionFactor
	for _, code := range resp.WarnFactors {
		desc := petmodel.GetFactorDescription(petmodel.FactorCode(code))
		warnFactors = append(warnFactors, dto.RestrictionFactor{
			Code:           code,
			Description:    desc.Description,
			SubDescription: desc.SubDescription,
		})
	}

	return dto.DonorApplication{
		ID:               resp.ID,
		RequestID:        resp.RequestID,
		DonorID:          resp.DonorID,
		DonorName:        resp.DonorName,
		DonorPhotos:      m.storage.BuildPhotoURLs(resp.DonorPhotos, *resp.UpdatedAt),
		DonorBloodGroup:  resp.DonorBloodGroup,
		Amount:           resp.Amount,
		WarnFactors:      warnFactors,
		CompensationType: resp.CompensationType,
		TaxiCompensation: resp.TaxiCompensation,
		Status:           string(resp.Status),
		IsConfirmed:      resp.IsConfirmed,
		CreatedAt:        resp.CreatedAt,
		UpdatedAt:        resp.UpdatedAt,
	}
}

func (m *BloodRequestMapper) ToResponseWithPotential(result *query.GetByPetIDResult) dto.BloodRequestDetail {
	if result == nil {
		return dto.BloodRequestDetail{}
	}

	detail := m.ToResponse(result.BloodRequest, &result.SuitableDonors)

	detail.PotentialDonors = make([]dto.DonorApplication, 0, len(result.PotentialDonors))
	for _, pd := range result.PotentialDonors {
		detail.PotentialDonors = append(detail.PotentialDonors, m.DonorResponseToApplication(pd))
	}

	return detail
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
				DonorPhotos:      m.storage.BuildPhotoURLs(app.DonorPhotos, *req.BloodRequest.UpdatedAt),
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
		ID:                       req.BloodRequest.ID,
		PetID:                    req.BloodRequest.PetID,
		BloodVolumeNeeded:        req.BloodRequest.BloodVolumeNeeded,
		BloodVolumeReserved:      req.BloodRequest.BloodVolumeReserved,
		BloodVolumeDonated:       req.BloodRequest.BloodVolumeDonated,
		Regions:                  req.BloodRequest.Regions,
		SmallPetsNotifyAllowed:   req.BloodRequest.SmallPetsNotifyAllowed,
		Description:              req.BloodRequest.AdvancedInfo.Description,
		PhotoURLs:                m.storage.BuildPhotoURLs(req.BloodRequest.AdvancedInfo.PhotoURLs, *req.BloodRequest.UpdatedAt),
		BloodGroupNames:          req.BloodRequest.BloodGroupNames,
		BloodComponentIDs:        req.BloodRequest.BloodComponentIDs,
		OnBoarding:               req.BloodRequest.OnBoarding,
		Status:                   string(req.BloodRequest.Status),
		PrioritySearch:           req.BloodRequest.PrioritySearch,
		IncludeUnknownBloodGroup: req.BloodRequest.IncludeUnknownBloodGroup,
		Responses:                donorApplications,
		AcceptedDonors:           acceptedDonorApplications,
		CompletedDonations:       completedDonations,
		SuitableDonors:           suitableDonorsCount,
		CreatedAt:                req.BloodRequest.CreatedAt,
		UpdatedAt:                req.BloodRequest.UpdatedAt,
		DeletedAt:                req.BloodRequest.DeletedAt,
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

// FromCreate converts a CreateBloodRequestBody DTO to a domain BloodRequest model.
func (m *BloodRequestMapper) FromCreate(body dto.CreateBloodRequestBody) *model.BloodRequest {
	req := &model.BloodRequest{
		PetID:                    body.PetID,
		BloodVolumeNeeded:        math.Round(body.BloodVolumeNeeded*10) / 10,
		BloodVolumeReserved:      0,
		Regions:                  body.Regions,
		SmallPetsNotifyAllowed:   body.SmallPetsNotifyAllowed,
		Status:                   model.BloodRequestStatusActive,
		BloodGroupNames:          body.BloodGroupNames,
		BloodComponentIDs:        body.BloodComponentIDs,
		OnBoarding:               []string{},
		PrioritySearch:           body.PrioritySearch,
		IncludeUnknownBloodGroup: body.IncludeUnknownBloodGroup,
		AdvancedInfo: model.AdvancedInfo{
			Description: body.Description,
			PhotoURLs:   []string{},
		},
	}

	return req
}
