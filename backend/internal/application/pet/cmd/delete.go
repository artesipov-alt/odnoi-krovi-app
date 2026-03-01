package cmd

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type DeleteHandler struct {
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.BloodRequestRepository
}

func NewDeleteHandler(petRepo pet.Repository, bloodReqRepo bloodsearch.BloodRequestRepository) *DeleteHandler {
	return &DeleteHandler{
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, petID string) error {
	exists, err := h.petRepo.ExistsByID(ctx, petID)
	if err != nil {
		return apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return apperrors.ErrPetNotFound
	}

	// Получаем все заявки на поиск крови, связанные с этим питомцем
	bloodRequests, err := h.bloodReqRepo.List(ctx, 0, 0, map[string]any{"pet_id": petID})
	if err != nil {
		return apperrors.Internal(err, "failed to list blood requests for pet")
	}

	// Удаляем каждую связанную заявку
	for _, req := range bloodRequests {
		if err := h.bloodReqRepo.Delete(ctx, req.ID); err != nil {
			slog.WarnContext(ctx, "Failed to delete blood request for pet", "blood_request_id", req.ID, "pet_id", petID, "error", err)
		}
	}

	if err := h.petRepo.Delete(ctx, petID); err != nil {
		return apperrors.Internal(err, "failed to delete pet")
	}

	return nil
}
