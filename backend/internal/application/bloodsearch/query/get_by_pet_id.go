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
	req, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, 0, err
	}

	if len(req.DonorApplications) != 0 {
		for _, app := range req.DonorApplications {
			donor, err := h.petRepo.GetByID(ctx, app.DonorID, pet.PetPreloadOptions{})
			if err != nil {
				return nil, 0, err
			}
			app.Amount = donor.CalculateDonationAmount()
		}
	}

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, req.BloodGroupNames)
	if err != nil {
		return nil, 0, err
	}

	return req, suitableDonors, nil
}
