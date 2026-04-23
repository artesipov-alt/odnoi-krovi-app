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
	bloodReqRepo  bloodsearch.BloodRequestRepository
	bonusRepo     bonus.Repository
	petService    *pet.PetService
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	donorRespRepo donor.Repository,
	bloodReqRepo bloodsearch.BloodRequestRepository,
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

	allPetsOpts := opts
	allPetsOpts.IgnoreSoftDelete = true
	allPets, err := h.petReadRepo.GetByUserID(ctx, userID, allPetsOpts)
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

	// Filter active pets for result
	activePets := make([]*model.Pet, 0, len(allPets))
	for _, pet := range allPets {
		if pet.DeletedAt == nil {
			activePets = append(activePets, pet)
		}
	}

	plannedDonations := make([]*donormodel.DonorResponse, 0, len(activePets))
	totalCompletedDonations := 0
	for _, pet := range allPets {
		applications := applicationsMap[pet.ID]
		var application *donormodel.DonorResponse
		for _, app := range applications {
			if app.IsActiveForDonation() {
				application = app
				break
			}
			if app.Status == donormodel.DonorResponseStatusCompleted && app.IsConfirmed {
				totalCompletedDonations++
			}
		}
		if application != nil && pet.DeletedAt == nil { // Only add planned for active pets
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
