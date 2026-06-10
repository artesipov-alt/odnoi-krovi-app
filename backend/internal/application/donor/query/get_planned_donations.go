package query

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type GetPlannedDonationsResult struct {
	ApplicationData    *donormodel.DonorResponse
	DonorPetData       *petmodel.Pet
	BloodSearchData    *bloodreqmodel.BloodRequestWithApplications
	RecipientPetData   *petmodel.Pet
	RecipientOwnerData *usermodel.User
	Bonuses            []*bonusmodel.Bonus
}

type PlannedDonationsHandler struct {
	donorRepo    donor.Repository
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
	userRepo     user.Repository
	bonusRepo    bonus.Repository
}

func NewPlannedDonationsHandler(donorRepo donor.Repository, petRepo pet.Repository, bloodReqRepo bloodsearch.BloodRequestRepository, userRepo user.Repository, bonusRepo bonus.Repository) *PlannedDonationsHandler {
	return &PlannedDonationsHandler{
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
		bonusRepo:    bonusRepo,
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
	applicationsMap, err := h.donorRepo.GetByPetIDs(ctx, petIDs, false)
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
			recipientOwner, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{
				WithIdentities: true,
			})
			if err != nil {
				return nil, apperrors.Internal(err, "failed to get recipient owner")
			}

			bonuses, err := h.bonusRepo.GetBonusesByDonorResponseID(ctx, application.ID)
			if err != nil {
				return nil, apperrors.Internal(err, "failed to get bonuses")
			}

			if len(bonuses) == 0 {
				bonuses = []*bonusmodel.Bonus{bonusmodel.NewLockBonus()}
			}

			result = append(result, &GetPlannedDonationsResult{
				ApplicationData:    application,
				BloodSearchData:    request,
				RecipientPetData:   recipientPet,
				DonorPetData:       dPet,
				RecipientOwnerData: recipientOwner,
				Bonuses:            bonuses,
			})
		}
	}

	return result, nil
}
