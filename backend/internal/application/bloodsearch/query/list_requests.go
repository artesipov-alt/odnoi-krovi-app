package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type ListRequestsHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
}

func NewListRequestsHandler(bloodRepo bloodsearch.BloodRequestRepository) *ListRequestsHandler {
	return &ListRequestsHandler{
		bloodRepo: bloodRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error) {
	requests, err := h.bloodRepo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	return requests, nil
}
