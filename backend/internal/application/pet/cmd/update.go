package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type UpdateHandler struct {
	petReadRepo  pet.PetReadRepository
	petWriteRepo pet.PetWriteRepository
	storage      filestorage.Repository
}

func NewUpdateHandler(petReadRepo pet.PetReadRepository, petWriteRepo pet.PetWriteRepository, storage filestorage.Repository) *UpdateHandler {
	return &UpdateHandler{
		petReadRepo:  petReadRepo,
		petWriteRepo: petWriteRepo,
		storage:      storage,
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

	// Нормализуем PhotoURLs - преобразуем полные URL обратно в относительные пути
	for i, url := range petInput.PhotoURLs {
		petInput.PhotoURLs[i] = h.storage.ExtractPathFromURL(url)
	}

	// Apply updates through aggregate method (controlled mutation)
	if err := existingPet.UpdateFrom(petInput); err != nil {
		return nil, apperrors.Internal(err, "failed to apply pet updates")
	}

	// Recalculate factors using aggregate method (encapsulates domain logic)
	existingPet.RecalculateFactors(time.Now())

	updatedPet, err := h.petWriteRepo.Update(ctx, id, existingPet)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update pet")
	}

	return updatedPet, nil
}
