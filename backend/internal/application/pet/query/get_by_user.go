package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// findActiveApplication finds the most recent active donor application from the list.
// Active means status is Accepted, Pending, or Completed but not confirmed.
func findActiveApplication(applications []*donormodel.DonorResponse) *donormodel.DonorResponse {
	for _, app := range applications {
		if app.Status == donormodel.DonorResponseStatusAccepted ||
			app.Status == donormodel.DonorResponseStatusPending ||
			(app.Status == donormodel.DonorResponseStatusCompleted && !app.IsConfirmed) {
			return app
		}
	}
	return nil
}

type GetByUserResult struct {
	Pets           []*model.Pet
	TotalPets      int
	TotalDonations int
}

type GetByUserHandler struct {
	petReadRepo   pet.PetReadRepository
	userRepo      user.Repository
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.BloodRequestRepository
	petService    *pet.PetService
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	donorRespRepo donor.Repository,
	bloodReqRepo bloodsearch.BloodRequestRepository,
	petService *pet.PetService,
) *GetByUserHandler {
	return &GetByUserHandler{
		petReadRepo:   petReadRepo,
		userRepo:      userRepo,
		bloodReqRepo:  bloodReqRepo,
		donorRespRepo: donorRespRepo,
		petService:    petService,
	}
}

func (h *GetByUserHandler) Handle(ctx context.Context, userID string, opts pet.PetPreloadOptions) (*GetByUserResult, error) {
	exists, err := h.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrUserNotFound
	}

	owner, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{
		WithDonorPreference: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get owner")
	}

	recoveryPeriodMonths := 0
	if owner.DonorPreference != nil {
		recoveryPeriodMonths = owner.DonorPreference.RecoveryPeriodMonths
	}

	pets, err := h.petReadRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	// Collect pet IDs for batch queries
	petIDs := make([]string, len(pets))
	for i, pet := range pets {
		petIDs[i] = pet.ID
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := h.donorRespRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodReqRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	plannedDonations := make([]*donormodel.DonorResponse, 0, len(pets))
	for _, pet := range pets {
		applications := applicationsMap[pet.ID]
		application := findActiveApplication(applications)
		if application != nil {
			if application.Status == donormodel.DonorResponseStatusAccepted ||
				(application.Status == donormodel.DonorResponseStatusCompleted && application.IsConfirmed == false) ||
				application.Status == donormodel.DonorResponseStatusPending {
				plannedDonations = append(plannedDonations, application)
			}
		}
		bloodReq := bloodReqsMap[pet.ID]

		h.petService.RecalculateFactorsAndStatus(pet, time.Now(), application, bloodReq)
		pet.RecoveryDays = h.petService.CalculateRecoveryDays(pet, recoveryPeriodMonths, time.Now())
	}

	return &GetByUserResult{
		Pets:           pets,
		TotalPets:      len(pets),
		TotalDonations: len(plannedDonations),
	}, nil
}
