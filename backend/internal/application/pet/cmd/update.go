package cmd

import (
	"context"

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
	// Load existing pet with all relations
	existingPet, err := h.petReadRepo.GetByID(ctx, id, pet.PetPreloadOptions{
		WithHealth:     true,
		WithTreatments: true,
		WithAnalyses:   true,
		WithBonuses:    true,
	})
	if err != nil {
		return nil, err // Доменная ошибка (например, ErrPetNotFound)
	}

	// Apply updates through aggregate method (controlled mutation)
	if err := existingPet.UpdateFrom(petInput); err != nil {
		return nil, apperrors.Internal(err, "failed to apply pet updates")
	}

	updatedPet, err := h.petWriteRepo.Update(ctx, id, existingPet)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update pet")
	}

	return updatedPet, nil
}
