package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetDonorByIDHandler struct {
	petReadRepo pet.PetReadRepository
	donorRepo   donor.Repository
}

func NewGetDonorByIDHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository) *GetDonorByIDHandler {
	return &GetDonorByIDHandler{
		petReadRepo: petReadRepo,
		donorRepo:   donorRepo,
	}
}

func (h *GetDonorByIDHandler) Handle(ctx context.Context, petID string, opts pet.PetPreloadOptions) (*petmodel.Pet, *donormodel.DonorResponse, error) {
	pet, err := h.petReadRepo.GetByID(ctx, petID, opts)
	if err != nil {
		return nil, nil, err
	}

	application, err := h.donorRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	pet.RecalculateFactors(now)
	pet.CalculateDonorStatus()

	return pet, application, nil
}
