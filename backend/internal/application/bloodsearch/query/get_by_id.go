package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
)

type GetByIDHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
}

func NewGetByIDHandler(bloodRepo bloodsearch.BloodRequestRepository) *GetByIDHandler {
	return &GetByIDHandler{
		bloodRepo: bloodRepo,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, id string) (*model.BloodRequest, error) {
	return h.bloodRepo.GetByID(ctx, id)
}
