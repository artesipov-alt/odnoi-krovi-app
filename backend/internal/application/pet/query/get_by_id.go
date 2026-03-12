package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetByIDHandler struct {
	petReadRepo  pet.PetReadRepository
	bloodReqRepo bloodsearch.BloodRequestRepository
}

func NewGetByIDHandler(petReadRepo pet.PetReadRepository, bloodReqRepo bloodsearch.BloodRequestRepository) *GetByIDHandler {
	return &GetByIDHandler{
		petReadRepo:  petReadRepo,
		bloodReqRepo: bloodReqRepo,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, petID string, opts pet.PetPreloadOptions) (*model.Pet, error) {
	p, err := h.petReadRepo.GetByID(ctx, petID, opts)
	if err != nil {
		return nil, err
	}

	p.RecalculateFactors(time.Now())
	p.CalculateDonorStatus()

	return p, nil
}
