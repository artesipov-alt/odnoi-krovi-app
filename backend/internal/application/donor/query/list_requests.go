package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type ListRequestsHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
}

func NewListRequestsHandler(bloodRepo bloodsearch.BloodRequestRepository, petRepo pet.Repository) *ListRequestsHandler {
	return &ListRequestsHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*donormodel.Recipient, error) {
	requests, err := h.bloodRepo.List(ctx, userID, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	return requests, nil
}
