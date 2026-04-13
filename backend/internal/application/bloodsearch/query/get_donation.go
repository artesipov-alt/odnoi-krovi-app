package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"

	bloodsearchmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetDonationHandler struct {
	petReadRepo pet.PetReadRepository
	donorRepo   donor.Repository
	bloodRepo   bloodsearch.BloodRequestRepository
	petService  *pet.PetService
}

type GetDonationResult struct {
	Application  *donormodel.DonorResponse
	BloodRequest *bloodsearchmodel.BloodRequestWithApplications
	DonorPet     *petmodel.Pet
	RecipientPet *petmodel.Pet
}

func NewGetDonationHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository, bloodRepo bloodsearch.BloodRequestRepository, petService *pet.PetService) *GetDonationHandler {
	return &GetDonationHandler{
		donorRepo:   donorRepo,
		bloodRepo:   bloodRepo,
		petReadRepo: petReadRepo,
		petService:  petService,
	}
}

func (h *GetDonationHandler) Handle(ctx context.Context, donorRespID string) (*GetDonationResult, error) {
	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorRespID)
	if err != nil {
		return nil, err
	}
	donorPet, err := h.petReadRepo.GetByID(ctx, application.DonorID, pet.PetPreloadOptions{
		WithHealth:     true,
		WithTreatments: true,
		WithAnalyses:   true,
	})
	if err != nil {
		return nil, err
	}
	donorBloodReq, err := h.bloodRepo.GetByPetID(ctx, donorPet.ID)
	if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
		return nil, err
	}

	h.petService.RecalculateFactorsAndStatus(donorPet, time.Now(), application, donorBloodReq)

	bloodReq, err := h.bloodRepo.GetByID(ctx, application.RequestID)
	if err != nil {
		return nil, err
	}

	bloodReq.RecalculateBloodAmount()
	bloodReq.RecalculateStatus()

	recipientPet, err := h.petReadRepo.GetByID(ctx, bloodReq.PetID, pet.PetPreloadOptions{
		WithHealth:     true,
		WithTreatments: true,
		WithAnalyses:   true,
	})
	if err != nil {
		return nil, err
	}
	h.petService.RecalculateFactorsAndStatus(recipientPet, time.Now(), nil, bloodReq)

	return &GetDonationResult{
		Application:  application,
		BloodRequest: bloodReq,
		DonorPet:     donorPet,
		RecipientPet: recipientPet,
	}, nil
}
