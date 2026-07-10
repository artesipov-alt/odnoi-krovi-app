package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type GetByPetIDHandler struct {
	bloodRepo    bloodsearch.Repository
	petRepo      pet.Repository
	bloodCounter *bloodsearch.BloodCounterService
}

func NewGetByPetIDHandler(
	bloodRepo bloodsearch.Repository,
	petRepo pet.Repository,
) *GetByPetIDHandler {
	return &GetByPetIDHandler{
		bloodRepo:    bloodRepo,
		petRepo:      petRepo,
		bloodCounter: bloodsearch.NewBloodCounterService(),
	}
}

func (h *GetByPetIDHandler) Handle(ctx context.Context, petID string) (*model.BloodRequestWithApplications, int, error) {
	bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
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
