package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type GetByIDHandler struct {
	bloodRepo    bloodsearch.Repository
	petRepo      pet.Repository
	bloodCounter *bloodsearch.BloodCounterService
}

func NewGetByIDHandler(bloodRepo bloodsearch.Repository, petRepo pet.Repository) *GetByIDHandler {
	return &GetByIDHandler{
		bloodRepo:    bloodRepo,
		petRepo:      petRepo,
		bloodCounter: bloodsearch.NewBloodCounterService(),
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, id string) (*model.BloodRequestWithApplications, int, error) {
	bloodReq, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)

	bloodReq.BloodRequest.SetBloodVolume(donated, reserved)

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, bloodReq.BloodRequest.BloodGroupNames)
	if err != nil {
		return nil, 0, err
	}

	return bloodReq, suitableDonors, nil
}
