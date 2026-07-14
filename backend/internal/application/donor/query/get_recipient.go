package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type RecipientDetailReadModel struct {
	Recipient       *bloodreqmodel.BloodRequestWithMatchingDonors
	AvilableBonuses []*bonusmodel.Bonus
}

type RecipientDetailHandler struct {
	donorRepo   donor.Repository
	petRepo     pet.Repository
	userRepo    user.Repository
	matchingSvc bloodsearch.MatchingService
	enricher    enrich.PetEnricher
	bonusSvc    *bonus.BonusService
}

func NewRecipientDetailHandler(donorRepo donor.Repository, petRepo pet.Repository, userRepo user.Repository, matchingSvc bloodsearch.MatchingService, enricher enrich.PetEnricher, bonusSvc *bonus.BonusService) *RecipientDetailHandler {
	return &RecipientDetailHandler{
		donorRepo:   donorRepo,
		petRepo:     petRepo,
		userRepo:    userRepo,
		matchingSvc: matchingSvc,
		enricher:    enricher,
		bonusSvc:    bonusSvc,
	}
}

func (h *RecipientDetailHandler) Handle(ctx context.Context, blodreqID string, userID string) (*RecipientDetailReadModel, error) {
	user, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{
		WithDonorPreference: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get user")
	}

	if user.DonorPreference == nil {
		return nil, apperrors.BadRequest("Настройки донора не заполнены")
	}

	preferredLocations := user.DonorPreference.PreferredLocationIDs

	recipient, err := h.donorRepo.GetRecipient(ctx, blodreqID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get recipient")
	}

	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets for user")
	}

	if _, err := h.enricher.RecalculateAll(ctx, pets, enrich.Options{RecoveryPeriodMonths: user.DonorPreference.RecoveryPeriodMonths}); err != nil {
		return nil, apperrors.Internal(err, "failed to recalculate all pets")
	}

	potentialDonors := petmodel.FilterDonors(pets)
	if len(potentialDonors) == 0 {
		return nil, apperrors.NotFound("потенциальные доноры не найдены")
	}

	for _, donorPet := range potentialDonors {
		if d, ok := h.matchingSvc.MatchDonor(recipient, donorPet, preferredLocations); ok {
			recipient.AddMatchingDonor(d)
		}
	}

	recipient.SetDefaultPrefs(user.DonorPreference.CompensationType, user.DonorPreference.TaxiCompensation)
	recipient.SyncPrivilegeAndPriority()

	arrears, err := h.bonusSvc.GetAggregatedBonuses(ctx, potentialDonors[0].Type, user.ID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get available bonuses")
	}

	return &RecipientDetailReadModel{
		Recipient:       recipient,
		AvilableBonuses: arrears,
	}, nil
}
