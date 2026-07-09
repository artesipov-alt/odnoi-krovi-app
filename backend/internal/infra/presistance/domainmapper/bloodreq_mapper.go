package domainmapper

import (
	"time"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// mapToRecipient maps ent.BloodSearchRequest to donormodel.Recipient
func RecipientToDomain(req *ent.BloodSearchRequest) *bloodreqmodel.BloodRequestWithMatchingDonors {
	if req == nil {
		return nil
	}

	recipient := &bloodreqmodel.BloodRequestWithMatchingDonors{
		BloodRequest: bloodreqmodel.BloodRequest{
			ID:                       req.ID,
			PetID:                    req.PetID,
			BloodVolumeNeeded:        req.BloodVolumeNeeded,
			PrioritySearch:           req.PrioritySearch,
			IncludeUnknownBloodGroup: req.IncludeUnknownBloodGroup,
			BloodGroupNames:          req.BloodGroupNames,
			SmallPetsNotifyAllowed:   req.SmallPetsNotifyAllowed,
			Regions:                  req.Regions,
			Status:                   bloodreqmodel.BloodRequestStatus(req.Status),
			CreatedAt:                &req.CreatedAt,
			UpdatedAt:                &req.UpdatedAt,
		},
	}

	if req.Edges.Pet != nil {
		recipient.RecipientData.PetName = req.Edges.Pet.Name
		recipient.RecipientData.OwnerID = req.Edges.Pet.UserID
		recipient.RecipientData.PetType = common.PetType(req.Edges.Pet.Type)
		recipient.RecipientData.PhotoURLs = req.Edges.Pet.PhotoUrls
		recipient.RecipientData.BloodGroupName = req.Edges.Pet.BloodGroup
		if req.Edges.Pet.Privilege != nil {
			recipient.RecipientData.Privilege = common.Privilege(*req.Edges.Pet.Privilege)
		}
	}

	return recipient
}

// BloodReqToDomain converts ENT BloodSearchRequest to domain BloodRequest
func BloodReqToDomain(entReq *ent.BloodSearchRequest) *bloodreqmodel.BloodRequestWithApplications {
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
				RejectedReason:   resp.RejectedReason,
				CreatedAt:        &resp.CreatedAt,
				UpdatedAt:        &resp.UpdatedAt,
			}
			if resp.Edges.Donor != nil {
				fullDonor := PetToDomain(resp.Edges.Donor)
				fullDonor.RecalculateFactors(now, false, false)
				app.DonorID = resp.Edges.Donor.ID
				app.DonorName = resp.Edges.Donor.Name
				app.DonorPhotos = resp.Edges.Donor.PhotoUrls
				app.DonorBloodGroup = fullDonor.BloodGroupName
				app.WarnFactors = fullDonor.WarnFactors
			}
			donorApps[i] = app
		}
	}

	return &bloodreqmodel.BloodRequestWithApplications{
		BloodRequest: bloodreqmodel.BloodRequest{
			ID:                       entReq.ID,
			PetID:                    entReq.PetID,
			BloodVolumeNeeded:        entReq.BloodVolumeNeeded,
			Regions:                  entReq.Regions,
			SmallPetsNotifyAllowed:   entReq.SmallPetsNotifyAllowed,
			Status:                   bloodreqmodel.BloodRequestStatus(entReq.Status),
			BloodGroupNames:          entReq.BloodGroupNames,
			BloodComponentIDs:        entReq.BloodComponentIds,
			OnBoarding:               entReq.OnBoarding,
			PrioritySearch:           entReq.PrioritySearch,
			IncludeUnknownBloodGroup: entReq.IncludeUnknownBloodGroup,
			AdvancedInfo: bloodreqmodel.AdvancedInfo{
				Description: entReq.Description,
				PhotoURLs:   entReq.PhotoUrls,
			},
			CreatedAt: &entReq.CreatedAt,
			UpdatedAt: &entReq.UpdatedAt,
			DeletedAt: entReq.DeletedAt,
		},
		DonorApplications: donorApps,
	}
}
