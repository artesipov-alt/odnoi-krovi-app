package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type DeleteHandler struct {
	petReadRepo  pet.PetReadRepository
	petWriteRepo pet.PetWriteRepository
	bloodReqRepo bloodsearch.Repository
}

func NewDeleteHandler(petReadRepo pet.PetReadRepository, petWriteRepo pet.PetWriteRepository, bloodReqRepo bloodsearch.Repository) *DeleteHandler {
	return &DeleteHandler{
		petReadRepo:  petReadRepo,
		petWriteRepo: petWriteRepo,
		bloodReqRepo: bloodReqRepo,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, petID string) error {
	exists, err := h.petReadRepo.Exists(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return apperrors.ErrPetNotFound
	}

	if err := h.petWriteRepo.DeleteWithRelations(ctx, petID); err != nil {
		return apperrors.Internal(err, "failed to delete pet")
	}

	return nil
}
