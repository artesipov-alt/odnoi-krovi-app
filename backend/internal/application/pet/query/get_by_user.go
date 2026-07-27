package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"

	bloodsearchmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	commonmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
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

	// Новое условие проверки потенциальных доноров.
	for _, p := range activePets {
		if !p.IsRecipient() {
			continue
		}
		potentialDonors, err := h.GetPotentialDonors(
			ctx,
			p.Type,
			fc.BloodReqs[p.ID].SearchingBloodGroupNames(),
			fc.BloodReqs[p.ID].SearchingRegions(),
			fc.BloodReqs[p.ID].BloodRequest.ID,
			owner.ID)
		if err != nil {
			return nil, err
		}
		if len(potentialDonors) == 0 {
			continue
		}
		if !IsCoversNededAmount(fc.BloodReqs[p.ID], potentialDonors) {
			continue
		}
		p.SetStatus(model.PetStatusBloodFound)
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

func (h *GetByUserHandler) GetPotentialDonors(ctx context.Context, petType commonmodel.PetType, bloodGroups, regions []string, callerUserID, requestID string) ([]*bloodsearchmodel.PotentialDonor, error) {
	potentialDonors, err := h.petReadRepo.FindPotentialDonors(ctx, pet.PotentialDonorsCriteria{
		PetType:          petType,
		BloodGroups:      bloodGroups,
		Regions:          regions,
		ExcludeRequestID: requestID,
		ExcludeOwnerID:   callerUserID,
	})
	if err != nil {
		return nil, err
	}

	// Извлекаем питомцев для батчевого Fetch и пересчёта статуса/факторов
	pets := make([]*model.Pet, len(potentialDonors))
	for i, pd := range potentialDonors {
		pets[i] = pd.Pet
	}

	// Батчевый Fetch для всех потенциальных доноров
	fc, err := h.enricher.Fetch(ctx, model.CollectIDs(pets))
	if err != nil {
		return nil, err
	}

	// Пересчёт статуса, факторов и recovery для каждого донора с его индивидуальным периодом
	donors := make([]*bloodsearchmodel.PotentialDonor, 0, len(potentialDonors))
	for _, pd := range potentialDonors {
		h.enricher.Recalculate(pd.Pet, fc, enrich.Options{
			RecoveryPeriodMonths: pd.RecoveryPeriodMonths,
		})

		if pd.Pet.PetStatus == model.PetStatusDonor {
			donors = append(donors, pd)
		}
	}

	return donors, nil
}

func IsCoversNededAmount(bloodReq *bloodsearchmodel.BloodRequestWithApplications, donorPets []*bloodsearchmodel.PotentialDonor) bool {
	for _, pd := range donorPets {
		if bloodReq.BloodRequest.IsCoversNededAmount(pd.Pet.CalculateDonationAmount()) {
			return true
		}
	}
	return false
}
