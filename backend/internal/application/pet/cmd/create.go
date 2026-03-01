package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type CreateHandler struct {
	petRepo  pet.PetWriteRepository
	userRepo user.Repository
}

func NewCreateHandler(petRepo pet.PetWriteRepository, userRepo user.Repository) *CreateHandler {
	return &CreateHandler{
		petRepo:  petRepo,
		userRepo: userRepo,
	}
}

func (h *CreateHandler) Handle(ctx context.Context, userID string, petInput *model.Pet) (*model.Pet, error) {
	// Проверяем, существует ли пользователь
	exists, err := h.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check user existence")
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	// Set the owner ID for the pet
	petInput.OwnerID = userID

	// Recalculate factors using aggregate method (encapsulates domain logic)
	petInput.RecalculateFactors(time.Now())

	newPet, err := h.petRepo.Create(ctx, petInput)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create pet")
	}

	return newPet, nil
}
