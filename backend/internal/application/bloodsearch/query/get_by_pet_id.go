package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type GetByPetIDHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
}

func NewGetByPetIDHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
) *GetByPetIDHandler {
	return &GetByPetIDHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
	}
}

func (h *GetByPetIDHandler) Handle(ctx context.Context, petID string) (*model.BloodRequest, int, error) {
	bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, 0, err
	}

	bloodReq.RecalculateBloodAmount()
	bloodReq.RecalculateStatus()

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, bloodReq.BloodGroupNames)
	if err != nil {
		return nil, 0, err
	}

	return bloodReq, suitableDonors, nil
}
