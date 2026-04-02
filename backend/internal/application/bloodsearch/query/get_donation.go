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
}

type GetDonationResult struct {
	Application  *donormodel.DonorResponse
	BloodRequest *bloodsearchmodel.BloodRequest
	DonorPet     *petmodel.Pet
	RecipientPet *petmodel.Pet
}

func NewGetDonationHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository, bloodRepo bloodsearch.BloodRequestRepository) *GetDonationHandler {
	return &GetDonationHandler{
		donorRepo:   donorRepo,
		bloodRepo:   bloodRepo,
		petReadRepo: petReadRepo,
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

	donorPet.RecalculateFactors(time.Now(), application, donorBloodReq)
	donorPet.CalculateStatus(application, donorBloodReq)

	bloodRequest, err := h.bloodRepo.GetByID(ctx, application.RequestID)
	if err != nil {
		return nil, err
	}

	bloodRequest.CalculateDonatedAmount()

	recipientPet, err := h.petReadRepo.GetByID(ctx, bloodRequest.PetID, pet.PetPreloadOptions{
		WithHealth:     true,
		WithTreatments: true,
		WithAnalyses:   true,
	})
	if err != nil {
		return nil, err
	}
	recipientPet.RecalculateFactors(time.Now(), nil, bloodRequest)
	recipientPet.CalculateStatus(nil, bloodRequest)

	return &GetDonationResult{
		Application:  application,
		BloodRequest: bloodRequest,
		DonorPet:     donorPet,
		RecipientPet: recipientPet,
	}, nil
}
