package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"

	bloodsearchmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type GetDonationHandler struct {
	petReadRepo  pet.PetReadRepository
	donorRepo    donor.Repository
	userRepo     user.Repository
	bloodRepo    bloodsearch.Repository
	bloodCounter *bloodsearch.BloodCounterService
}

type GetDonationResult struct {
	Application    *donormodel.DonorResponse
	BloodRequest   *bloodsearchmodel.BloodRequestWithApplications
	DonorPet       *petmodel.Pet
	DonorOwnerData *usermodel.User
	RecipientPet   *petmodel.Pet
}

func NewGetDonationHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository, userRepo user.Repository, bloodRepo bloodsearch.Repository) *GetDonationHandler {
	return &GetDonationHandler{
		donorRepo:    donorRepo,
		bloodRepo:    bloodRepo,
		petReadRepo:  petReadRepo,
		userRepo:     userRepo,
		bloodCounter: bloodsearch.NewBloodCounterService(),
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
	donorOwnerData, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return nil, err
	}

	donorBloodReq, err := h.bloodRepo.GetByPetID(ctx, donorPet.ID)
	if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
		return nil, err
	}

	donorPet.RecalculateStatus(time.Now(), pet.BuildDonationContext(application, donorBloodReq))

	bloodReq, err := h.bloodRepo.GetByID(ctx, application.RequestID)
	if err != nil {
		return nil, err
	}

	donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)

	bloodReq.BloodRequest.SetBloodVolume(donated, reserved)
	bloodReq.RecalculateStatus()

	recipientPet, err := h.petReadRepo.GetByID(ctx, bloodReq.BloodRequest.PetID, pet.PetPreloadOptions{
		WithHealth:     true,
		WithTreatments: true,
		WithAnalyses:   true,
	})
	if err != nil {
		return nil, err
	}
	recipientPet.RecalculateStatus(time.Now(), pet.BuildDonationContext(nil, bloodReq))

	return &GetDonationResult{
		Application:    application,
		BloodRequest:   bloodReq,
		DonorPet:       donorPet,
		DonorOwnerData: donorOwnerData,
		RecipientPet:   recipientPet,
	}, nil
}
