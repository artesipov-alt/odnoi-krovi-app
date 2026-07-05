package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type GetByIDHandler struct {
	bloodRepo bloodsearch.Repository
	petRepo   pet.Repository
}

func NewGetByIDHandler(bloodRepo bloodsearch.Repository, petRepo pet.Repository) *GetByIDHandler {
	return &GetByIDHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, id string) (*model.BloodRequestWithApplications, int, error) {
	bloodReq, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	bloodReq.RecalculateBloodAmount()

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, bloodReq.BloodGroupNames)
	if err != nil {
		return nil, 0, err
	}

	return bloodReq, suitableDonors, nil
}
