package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type RecipientDetailReadModel struct {
	Recipient       *bloodreqmodel.BloodRequestWithMatchingDonors
	AvilableBonuses []*bonusmodel.Bonus
}

type RecipientDetailHandler struct {
	donorRepo    donor.Repository
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.Repository
	userRepo     user.Repository
	matchingSvc  bloodsearch.MatchingService
	bonusSvc     *bonus.BonusService
}

func NewRecipientDetailHandler(donorRepo donor.Repository, petRepo pet.Repository, bloodReqRepo bloodsearch.Repository, userRepo user.Repository, matchingSvc bloodsearch.MatchingService, bonusSvc *bonus.BonusService) *RecipientDetailHandler {
	return &RecipientDetailHandler{
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
		matchingSvc:  matchingSvc,
		bonusSvc:     bonusSvc,
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

	// Collect pet IDs for batch queries
	petIDs := make([]string, len(pets))
	for i, pet := range pets {
		petIDs[i] = pet.ID
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := h.donorRepo.GetByPetIDs(ctx, petIDs, false)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodReqRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	for _, potentialDonorPet := range pets {
		applications := applicationsMap[potentialDonorPet.ID]
		var application *donormodel.DonorResponse
		for _, app := range applications {
			if app.IsActiveForDonation() {
				application = app
				break
			}
		}
		bloodReq := bloodReqsMap[potentialDonorPet.ID]
		potentialDonorPet.RecalculateStatus(time.Now(), pet.BuildDonationContext(application, bloodReq))
	}

	potentialDonors := petmodel.FilterDonors(pets)

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
