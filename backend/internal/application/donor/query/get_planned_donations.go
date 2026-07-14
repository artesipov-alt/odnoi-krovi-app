package query

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
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
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.Repository
	userRepo     user.Repository
	bonusRepo    bonus.Repository
	enricher     enrich.PetEnricher
}

func NewPlannedDonationsHandler(petRepo pet.Repository, bloodReqRepo bloodsearch.Repository, userRepo user.Repository, bonusRepo bonus.Repository, enricher enrich.PetEnricher) *PlannedDonationsHandler {
	return &PlannedDonationsHandler{
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
		bonusRepo:    bonusRepo,
		enricher:     enricher,
	}
}

func (h *PlannedDonationsHandler) Handle(ctx context.Context, userID string) ([]*GetPlannedDonationsResult, error) {
	donorPets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{WithAll: true})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets for user")
	}

	petIDs := make([]string, len(donorPets))
	for i, dPet := range donorPets {
		petIDs[i] = dPet.ID
	}

	fc, err := h.enricher.Fetch(ctx, petIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*GetPlannedDonationsResult, 0, len(donorPets))
	for _, dPet := range donorPets {
		application := h.enricher.Recalculate(dPet, fc, enrich.Options{})
		if dPet.PetStatus != petmodel.PetStatusPlannedDonation {
			continue
		}

		enriched, err := h.enrichPlannedDonation(ctx, dPet, application)
		if err != nil {
			if errors.Is(err, apperrors.ErrBloodRequestNotFound) {
				continue
			}
			return nil, err
		}
		result = append(result, enriched)
	}

	return result, nil
}

// enrichPlannedDonation достраивает результат для одного активного отклика:
// заявку, питомца-реципиента, его владельца и бонусы.
func (h *PlannedDonationsHandler) enrichPlannedDonation(ctx context.Context, donorPet *petmodel.Pet, application *donormodel.DonorResponse) (*GetPlannedDonationsResult, error) {
	request, err := h.bloodReqRepo.GetByApplicationID(ctx, application.ID, false)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood request")
	}

	recipientPet, err := h.petRepo.GetByID(ctx, request.BloodRequest.PetID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pet")
	}

	recipientOwner, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{WithIdentities: true})
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

	return &GetPlannedDonationsResult{
		ApplicationData:    application,
		BloodSearchData:    request,
		RecipientPetData:   recipientPet,
		DonorPetData:       donorPet,
		RecipientOwnerData: recipientOwner,
		Bonuses:            bonuses,
	}, nil
}
