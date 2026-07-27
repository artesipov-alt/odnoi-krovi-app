package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type GetByIDHandler struct {
	petReadRepo pet.PetReadRepository
	userRepo    user.Repository
	enricher    enrich.PetEnricher
}

func NewGetByIDHandler(petReadRepo pet.PetReadRepository, userRepo user.Repository, enricher enrich.PetEnricher) *GetByIDHandler {
	return &GetByIDHandler{
		petReadRepo: petReadRepo,
		userRepo:    userRepo,
		enricher:    enricher,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, petID string, opts pet.PetPreloadOptions) (*model.Pet, error) {
	p, err := h.petReadRepo.GetByID(ctx, petID, opts)
	if err != nil {
		return nil, err
	}

	owner, err := h.userRepo.GetByID(ctx, p.OwnerID, user.UserPreloadOptions{
		WithDonorPreference: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get owner")
	}

	recoveryPeriodMonths := 0
	if owner.DonorPreference != nil {
		recoveryPeriodMonths = owner.DonorPreference.RecoveryPeriodMonths
	}

	fc, err := h.enricher.Fetch(ctx, []string{petID})
	if err != nil {
		return nil, err
	}
	h.enricher.Recalculate(p, fc, enrich.Options{RecoveryPeriodMonths: recoveryPeriodMonths})

	return p, nil
}
