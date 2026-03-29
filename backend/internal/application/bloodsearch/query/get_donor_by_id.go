package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetDonorByIDHandler struct {
	petReadRepo pet.PetReadRepository
	donorRepo   donor.Repository
	bloodRepo   bloodsearch.BloodRequestRepository
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
	if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
		return nil, nil, err
	}

	bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
		return nil, nil, err
	}

	pet.RecalculateFactors(time.Now(), application, bloodReq)
	pet.CalculateStatus(application, bloodReq)

	return pet, application, nil
}
