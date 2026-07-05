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
	bloodRepo   bloodsearch.Repository
	petService  pet.PetService
}

func NewGetDonorByIDHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository, bloodRepo bloodsearch.Repository, petService pet.PetService) *GetDonorByIDHandler {
	return &GetDonorByIDHandler{
		petReadRepo: petReadRepo,
		donorRepo:   donorRepo,
		bloodRepo:   bloodRepo,
		petService:  petService,
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
	// if application == nil {
	// 	return nil, nil, apperrors.ErrDonorResponseNotFound
	// }

	bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
		return nil, nil, err
	}

	if bloodReq != nil {
		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()
	}

	h.petService.RecalculateFactorsAndStatus(pet, time.Now(), application, bloodReq)

	return pet, application, nil
}
