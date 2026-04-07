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

	result := make([]*GetPlannedDonationsResult, len(donorPets))
	for _, dPet := range donorPets {
		application, err := h.donorRepo.GetByPetID(ctx, dPet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
			return nil, apperrors.Internal(err, "failed to get donor application")
		}
		if application != nil {
			request, err := h.bloodReqRepo.GetByApplicationID(ctx, application.ID)
			if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
				return nil, apperrors.Internal(err, "failed to get blood request")
			}
			recipientPet, err := h.petRepo.GetByID(ctx, request.PetID, pet.PetPreloadOptions{})
			if err != nil {
				return nil, apperrors.Internal(err, "failed to get pet")
			}
			result = append(result, &GetPlannedDonationsResult{
				ApplicationData:  *application,
				BloodSearchData:  *request,
				RecipientPetData: *recipientPet,
			})
		}
	}

	return result, nil
}
