package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
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
	petReadRepo pet.PetReadRepository
	userRepo    user.Repository
	bonusRepo   bonus.Repository
	enricher    enrich.PetEnricher
}

func NewGetByUserHandler(
	petReadRepo pet.PetReadRepository,
	userRepo user.Repository,
	bonusRepo bonus.Repository,
	enricher enrich.PetEnricher,
) *GetByUserHandler {
	return &GetByUserHandler{
		petReadRepo: petReadRepo,
		userRepo:    userRepo,
		bonusRepo:   bonusRepo,
		enricher:    enricher,
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

	// Без opts.SetIgnoreSoftDelete() — репозиторий сам вернёт только активных питомцев.
	// Soft-deleted больше не нужны здесь: TotalCompletedDonations считается отдельным
	// SQL-агрегатом (CountFullyCompletedDonations), который сам учитывает удалённых.
	activePets, err := h.petReadRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	petIDs := make([]string, len(activePets))
	for i, p := range activePets {
		petIDs[i] = p.ID
	}

	fc, err := h.enricher.Fetch(ctx, petIDs)
	if err != nil {
		return nil, err
	}

	totalPlannedDonations := 0
	for _, p := range activePets {
		app := h.enricher.Recalculate(p, fc, enrich.Options{RecoveryPeriodMonths: recoveryPeriodMonths})
		if app != nil && app.IsActiveForDonation() {
			totalPlannedDonations++
		}
	}

	totalCompletedDonations, err := h.enricher.CountFullyCompletedDonations(ctx, userID)
	if err != nil {
		return nil, err
	}

	assignedBonuses, err := h.bonusRepo.GetAssignedBonuses(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get assigned bonuses")
	}

	return &GetByUserResult{
		Pets:                    activePets,
		TotalPets:               len(activePets),
		TotalPlannedDonations:   totalPlannedDonations,
		TotalCompletedDonations: totalCompletedDonations,
		TotalPrioritySearch:     owner.PrioritySearchCount,
		TotalBonuses:            len(assignedBonuses),
	}, nil
}
