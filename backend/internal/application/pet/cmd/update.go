package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type UpdateHandler struct {
	petReadRepo  pet.PetReadRepository
	petWriteRepo pet.PetWriteRepository
}

func NewUpdateHandler(petReadRepo pet.PetReadRepository, petWriteRepo pet.PetWriteRepository) *UpdateHandler {
	return &UpdateHandler{
		petReadRepo:  petReadRepo,
		petWriteRepo: petWriteRepo,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, id string, petInput *model.Pet) (*model.Pet, error) {
	exists, err := h.petReadRepo.Exists(ctx, id)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return nil, apperrors.ErrPetNotFound
	}

	// Calculate static stop factors (stored in DB) and warn factors
	stopFactors := petInput.GetStaticStopFactors()
	petInput.StopFactors = make([]string, len(stopFactors))
	for i, f := range stopFactors {
		petInput.StopFactors[i] = string(f)
	}

	warnFactors := petInput.GetWarnFactors(time.Now())
	petInput.WarnFactors = make([]string, len(warnFactors))
	for i, f := range warnFactors {
		petInput.WarnFactors[i] = string(f)
	}

	updatedPet, err := h.petWriteRepo.Update(ctx, id, petInput)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update pet")
	}

	return updatedPet, nil
}
