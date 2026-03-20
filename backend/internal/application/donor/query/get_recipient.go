package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type RecipientDetailHandler struct {
	donorRepo donor.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
}

func NewRecipientDetailHandler(donorRepo donor.Repository, petRepo pet.Repository, userRepo user.Repository) *RecipientDetailHandler {
	return &RecipientDetailHandler{
		donorRepo: donorRepo,
		petRepo:   petRepo,
		userRepo:  userRepo,
	}
}

func (h *RecipientDetailHandler) Handle(ctx context.Context, blodreqID string, userID string) (*donormodel.Recipient, error) {
	recipient, err := h.donorRepo.GetRecipient(ctx, blodreqID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get recipient")
	}

	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll:      true,
		WithBloodReq: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets for user")
	}

	for i, _ := range pets {
		pets[i].RecalculateFactors(time.Now())
		pets[i].CalculateDonorStatus()
	}

	potentialDonors := petmodel.FilterDonors(pets)

	for _, donorPet := range potentialDonors {
		recipient.MatchDonor(donorPet)
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
