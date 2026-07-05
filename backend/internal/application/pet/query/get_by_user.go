package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type GetByUserResult struct {
	Pets                    []*model.Pet
	TotalPets               int
	TotalPlannedDonations   int
	TotalCompletedDonations int
	TotalPrioritySearch     int
	TotalBonuses            int
}

type GetByUserHandler struct {
	petReadRepo   pet.PetReadRepository
	userRepo      user.Repository
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.Repository
	bonusRepo     bonus.Repository
	petService    *pet.PetService
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	donorRespRepo donor.Repository,
	bloodReqRepo bloodsearch.Repository,
	bonusRepo bonus.Repository,
	petService *pet.PetService,
) *GetByUserHandler {
	return &GetByUserHandler{
		petReadRepo:   petReadRepo,
		userRepo:      userRepo,
		bloodReqRepo:  bloodReqRepo,
		donorRespRepo: donorRespRepo,
		bonusRepo:     bonusRepo,
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

	opts.SetIgnoreSoftDelete()
	allPets, err := h.petReadRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	// Collect pet IDs for batch queries (all pets, including deleted)
	petIDs := make([]string, len(allPets))
	for i, pet := range allPets {
		petIDs[i] = pet.ID
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := h.donorRespRepo.GetByPetIDs(ctx, petIDs, true)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodReqRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	// ---- Stage 1: count completed donations for all pets (including deleted) ----
	totalCompletedDonations := 0
	for _, pet := range allPets {
		for _, app := range applicationsMap[pet.ID] {
			if app.IsCompleted() {
				totalCompletedDonations++
			}
		}
	}

	// ---- Stage 2: process only active pets ----
	activePets := make([]*model.Pet, 0, len(allPets))
	plannedDonations := make([]*donormodel.DonorResponse, 0)

	for _, pet := range allPets {
		if pet.IsDeleted() {
			continue
		}
		activePets = append(activePets, pet)

		application := findActiveDonation(applicationsMap[pet.ID])
		if application != nil {
			plannedDonations = append(plannedDonations, application)
		}

		bloodReq := bloodReqsMap[pet.ID]
		h.petService.RecalculateFactorsAndStatus(pet, time.Now(), application, bloodReq)
		pet.RecoveryDays = h.petService.CalculateRecoveryDays(pet, recoveryPeriodMonths, time.Now())
	}

	assignedBonuses, err := h.bonusRepo.GetAssignedBonuses(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get assigned bonuses")
	}

	return &GetByUserResult{
		Pets:                    activePets,
		TotalPets:               len(activePets),
		TotalPlannedDonations:   len(plannedDonations),
		TotalCompletedDonations: totalCompletedDonations,
		TotalPrioritySearch:     owner.PrioritySearchCount,
		TotalBonuses:            len(assignedBonuses),
	}, nil
}

func findActiveDonation(applications []*donormodel.DonorResponse) *donormodel.DonorResponse {
	for _, app := range applications {
		if app.IsActiveForDonation() {
			return app
		}
	}
	return nil
}
