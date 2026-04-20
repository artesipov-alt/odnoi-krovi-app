package query

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type GetPlannedDonationsResult struct {
	ApplicationData  donormodel.DonorResponse
	DonorPetData     petmodel.Pet
	BloodSearchData  bloodreqmodel.BloodRequestWithApplications
	RecipientPetData petmodel.Pet
}

type PlannedDonationsHandler struct {
	donorRepo    donor.Repository
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
	userRepo     user.Repository
}

func NewPlannedDonationsHandler(donorRepo donor.Repository, petRepo pet.Repository, bloodReqRepo bloodsearch.BloodRequestRepository, userRepo user.Repository) *PlannedDonationsHandler {
	return &PlannedDonationsHandler{
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
	}
}

func (h *PlannedDonationsHandler) Handle(ctx context.Context, userID string) ([]*GetPlannedDonationsResult, error) {
	donorPets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets for user")
	}

	// Collect pet IDs for batch queries
	petIDs := make([]string, len(donorPets))
	for i, dPet := range donorPets {
		petIDs[i] = dPet.ID
	}

	// Batch fetch applications
	applicationsMap, err := h.donorRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	result := make([]*GetPlannedDonationsResult, 0, len(donorPets))
	for _, dPet := range donorPets {
		applications := applicationsMap[dPet.ID]
		var application *donormodel.DonorResponse
		for _, app := range applications {
			if app.IsActiveForDonation() {
				application = app
				break
			}
		}
		if application != nil {
			request, err := h.bloodReqRepo.GetByApplicationID(ctx, application.ID, false)
			if err != nil {
				if !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
					return nil, apperrors.Internal(err, "failed to get blood request")
				}
				// Skip if blood request not found
				continue
			}
			recipientPet, err := h.petRepo.GetByID(ctx, request.PetID, pet.PetPreloadOptions{})
			if err != nil {
				return nil, apperrors.Internal(err, "failed to get pet")
			}
			result = append(result, &GetPlannedDonationsResult{
				ApplicationData:  *application,
				BloodSearchData:  *request,
				RecipientPetData: *recipientPet,
				DonorPetData:     *dPet,
			})
		}
	}

	return result, nil
}
