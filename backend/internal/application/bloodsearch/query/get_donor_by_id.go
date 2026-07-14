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
	petReadRepo  pet.PetReadRepository
	donorRepo    donor.Repository
	bloodRepo    bloodsearch.Repository
	bloodCounter *bloodsearch.BloodCounterService
}

func NewGetDonorByIDHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository, bloodRepo bloodsearch.Repository) *GetDonorByIDHandler {
	return &GetDonorByIDHandler{
		petReadRepo:  petReadRepo,
		donorRepo:    donorRepo,
		bloodRepo:    bloodRepo,
		bloodCounter: bloodsearch.NewBloodCounterService(),
	}
}

func (h *GetDonorByIDHandler) Handle(ctx context.Context, petID string, opts pet.PetPreloadOptions) (*petmodel.Pet, *donormodel.DonorResponse, error) {
	donorPet, err := h.petReadRepo.GetByID(ctx, petID, opts)
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
		donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)

		bloodReq.BloodRequest.SetBloodVolume(donated, reserved)
		bloodReq.RecalculateStatus()
	}

	donorPet.RecalculateStatus(time.Now(), pet.BuildDonationContext(application, bloodReq))

	return donorPet, application, nil
}
