package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type RecipientDetailHandler struct {
	donorRepo    donor.Repository
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
	userRepo     user.Repository
	matchingSvc  bloodsearch.MatchingService
}

func NewRecipientDetailHandler(donorRepo donor.Repository, petRepo pet.Repository, bloodReqRepo bloodsearch.BloodRequestRepository, userRepo user.Repository, matchingSvc bloodsearch.MatchingService) *RecipientDetailHandler {
	return &RecipientDetailHandler{
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
		matchingSvc:  matchingSvc,
	}
}

func (h *RecipientDetailHandler) Handle(ctx context.Context, blodreqID string, userID string) (*bloodreqmodel.BloodRequestWithMatchingDonors, error) {
	recipient, err := h.donorRepo.GetRecipient(ctx, blodreqID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get recipient")
	}

	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets for user")
	}

	for _, pet := range pets {
		application, err := h.donorRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
			return nil, apperrors.Internal(err, "failed to get donor application")
		}
		bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			return nil, apperrors.Internal(err, "failed to get blood request")
		}
		pet.RecalculateFactors(time.Now(), application, bloodReq)
		pet.CalculateStatus(application, bloodReq)
	}

	potentialDonors := petmodel.FilterDonors(pets)

	for _, donorPet := range potentialDonors {
		h.matchingSvc.MatchDonor(recipient, donorPet)
	}

	user, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{
		WithDonorPreference: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get user for recipient details")
	}

	recipient.SetDefaultPrefs(user.DonorPreference.CompensationType, user.DonorPreference.TaxiCompensation)

	return recipient, nil
}
