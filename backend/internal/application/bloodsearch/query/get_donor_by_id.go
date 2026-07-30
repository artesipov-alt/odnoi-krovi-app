package query

import (
	"context"
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type GetDonorByIDResult struct {
	Donor       *petmodel.Pet
	Application *donormodel.DonorResponse
}

type GetDonorByIDHandler struct {
	petReadRepo  pet.PetReadRepository
	donorRepo    donor.Repository
	userRepo     user.Repository
	bloodRepo    bloodsearch.Repository
	bloodCounter *bloodsearch.BloodCounterService
	enricher     enrich.PetEnricher
}

func NewGetDonorByIDHandler(petReadRepo pet.PetReadRepository, donorRepo donor.Repository, userRepo user.Repository, bloodRepo bloodsearch.Repository, enricher enrich.PetEnricher) *GetDonorByIDHandler {
	return &GetDonorByIDHandler{
		petReadRepo:  petReadRepo,
		donorRepo:    donorRepo,
		userRepo:     userRepo,
		bloodRepo:    bloodRepo,
		bloodCounter: bloodsearch.NewBloodCounterService(),
		enricher:     enricher,
	}
}

func (h *GetDonorByIDHandler) Handle(ctx context.Context, petID string, opts pet.PetPreloadOptions) (*petmodel.Pet, *donormodel.DonorResponse, error) {
	donorPet, err := h.petReadRepo.GetByID(ctx, petID, opts)
	if err != nil {
		return nil, nil, err
	}

	fc, err := h.enricher.Fetch(ctx, []string{petID})
	if err != nil {
		return nil, nil, err
	}

	application, err := h.donorRepo.GetByPetID(ctx, petID)
	if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
		return nil, nil, err
	}

	h.enricher.Recalculate(donorPet, fc, enrich.Options{})

	// Если отклика нет, инициализируем пустой с компенсацией из предпочтений владельца
	if application == nil {
		prefs := usermodel.DefaultDonorPreference()

		owner, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{
			WithDonorPreference: true,
		})
		if err == nil && owner.DonorPreference != nil {
			prefs = owner.DonorPreference
		}

		application = &donormodel.DonorResponse{
			DonorPrefs: donormodel.DonorPrefs{
				PreferredLocationIDs: prefs.PreferredLocationIDs,
				CompensationType:     string(prefs.CompensationType),
				TaxiCompensation:     prefs.TaxiCompensation,
			},
			Status: donormodel.DonorResponseStatusPending,
		}
	}

	// if bloodReq := fc.BloodReqs[petID]; bloodReq != nil {
	// 	donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)

	// 	bloodReq.BloodRequest.SetBloodVolume(donated, reserved)
	// 	bloodReq.RecalculateStatus()
	// }

	return donorPet, application, nil
}
