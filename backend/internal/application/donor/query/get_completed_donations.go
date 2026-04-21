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

type GetCompletedDonationsResult struct {
	ApplicationData  donormodel.DonorResponse
	DonorPetData     petmodel.Pet
	BloodSearchData  bloodreqmodel.BloodRequestWithApplications
	RecipientPetData petmodel.Pet
}

type CompletedDonationsHandler struct {
	donorRepo    donor.Repository
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
	userRepo     user.Repository
}

func NewCompletedDonationsHandler(donorRepo donor.Repository, petRepo pet.Repository, bloodReqRepo bloodsearch.BloodRequestRepository, userRepo user.Repository) *CompletedDonationsHandler {
	return &CompletedDonationsHandler{
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
	}
}

func (h *CompletedDonationsHandler) Handle(ctx context.Context, userID string) ([]*GetCompletedDonationsResult, error) {
	donorPets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll:          true,
		IgnoreSoftDelete: true,
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

	result := make([]*GetCompletedDonationsResult, 0, len(donorPets))
	for _, dPet := range donorPets {
		applications := applicationsMap[dPet.ID]
		for _, app := range applications {
			if app.IsClosedForDonation() {
				// TODO: N+1 На каждую заявку тянется по одному запросу. Нужно сделать общий метод.
				request, err := h.bloodReqRepo.GetByApplicationID(ctx, app.ID, true)
				if err != nil {
					if !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
						return nil, apperrors.Internal(err, "failed to get blood request")
					}
					// Skip if blood request not found
					continue
				}
				recipientPet, err := h.petRepo.GetByID(ctx, request.PetID, pet.PetPreloadOptions{IgnoreSoftDelete: true})
				if err != nil {
					return nil, apperrors.Internal(err, "failed to get pet")
				}
				result = append(result, &GetCompletedDonationsResult{
					ApplicationData:  *app,
					BloodSearchData:  *request,
					RecipientPetData: *recipientPet,
					DonorPetData:     *dPet,
				})
			}
		}
	}

	return result, nil
}
