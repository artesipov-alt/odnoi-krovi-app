package domainmapper

import (
	"time"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	recipientmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/recipient/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// mapToRecipient maps ent.BloodSearchRequest to donormodel.Recipient
func RecipientToDomain(req *ent.BloodSearchRequest) *recipientmodel.Recipient {
	if req == nil {
		return nil
	}

	recipient := &recipientmodel.Recipient{
		ID:                       req.ID,
		PetID:                    req.PetID,
		BloodVolumeRemaining:     req.BloodVolumeNeeded - req.BloodVolumeReserved,
		PrioritySearch:           req.PrioritySearch,
		IncludeUnknownBloodGroup: req.IncludeUnknownBloodGroup,
		SearchingBloodNames:      req.BloodGroupNames,
		SmallPetsNotifyAllowed:   req.SmallPetsNotifyAllowed,
		Status:                   string(req.Status),
	}

	if req.Edges.Pet != nil {
		recipient.PetName = req.Edges.Pet.Name
		recipient.PetType = petmodel.PetType(req.Edges.Pet.Type)
		recipient.PhotoURLs = req.Edges.Pet.PhotoUrls
		if req.Edges.Pet.Edges.BloodGroupRef != nil {
			recipient.BloodGroupName = req.Edges.Pet.Edges.BloodGroupRef.BloodGroup
		}
	}

	return recipient
}

// BloodReqToDomain converts ENT BloodSearchRequest to domain BloodRequest
func BloodReqToDomain(entReq *ent.BloodSearchRequest) *bloodreqmodel.BloodRequest {
	if entReq == nil {
		return nil
	}
	now := time.Now()
	// Map responses to DonorApplications
	var donorApps []donormodel.DonorResponse
	if entReq.Edges.Responses != nil {
		donorApps = make([]donormodel.DonorResponse, len(entReq.Edges.Responses))
		for i, resp := range entReq.Edges.Responses {
			app := donormodel.DonorResponse{
				ID:               resp.ID,
				RequestID:        entReq.ID,
				Amount:           resp.Amount,
				CompensationType: string(resp.CompensationType),
				TaxiCompensation: resp.TaxiCompensation,
				Status:           donormodel.DonorResponseStatus(resp.Status),
				IsConfirmed:      resp.IsConfirmed,
			}
			if resp.Edges.Donor != nil {
				fullDonor := PetToDomain(resp.Edges.Donor)
				fullDonor.RecalculateFactors(now, &app, nil)
				app.DonorID = resp.Edges.Donor.ID
				app.DonorName = resp.Edges.Donor.Name
				app.DonorPhotos = resp.Edges.Donor.PhotoUrls
				app.DonorBloodGroup = resp.Edges.Donor.Edges.BloodGroupRef.BloodGroup
				app.WarnFactors = fullDonor.WarnFactors
			}
			donorApps[i] = app
		}
	}

	return &bloodreqmodel.BloodRequest{
		ID:                       entReq.ID,
		PetID:                    entReq.PetID,
		BloodVolumeNeeded:        entReq.BloodVolumeNeeded,
		BloodVolumeReserved:      entReq.BloodVolumeReserved,
		Regions:                  entReq.Regions,
		SmallPetsNotifyAllowed:   entReq.SmallPetsNotifyAllowed,
		Status:                   bloodreqmodel.BloodRequestStatus(entReq.Status),
		Description:              entReq.Description,
		PhotoURLs:                entReq.PhotoUrls,
		BloodGroupNames:          entReq.BloodGroupNames,
		BloodComponentIDs:        entReq.BloodComponentIds,
		OnBoarding:               entReq.OnBoarding,
		DonorApplications:        donorApps,
		PrioritySearch:           entReq.PrioritySearch,
		IncludeUnknownBloodGroup: entReq.IncludeUnknownBloodGroup,
		CreatedAt:                &entReq.CreatedAt,
		UpdatedAt:                &entReq.UpdatedAt,
		DeletedAt:                entReq.DeletedAt,
	}
}
