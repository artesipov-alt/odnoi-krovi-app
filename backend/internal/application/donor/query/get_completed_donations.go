package query

import (
	"context"
	"errors"

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

type GetCompletedDonationsResult struct {
	ApplicationData  donormodel.DonorResponse
	DonorPetData     petmodel.Pet
	BloodSearchData  bloodreqmodel.BloodRequestWithApplications
	RecipientPetData petmodel.Pet
	Bonuses          []*bonusmodel.Bonus
}

type CompletedDonationsHandler struct {
	donorRepo    donor.Repository
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.Repository
	userRepo     user.Repository
	bonusRepo    bonus.Repository
}

func NewCompletedDonationsHandler(donorRepo donor.Repository, petRepo pet.Repository, bloodReqRepo bloodsearch.Repository, userRepo user.Repository, bonusRepo bonus.Repository) *CompletedDonationsHandler {
	return &CompletedDonationsHandler{
		donorRepo:    donorRepo,
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		userRepo:     userRepo,
		bonusRepo:    bonusRepo,
	}
}

func (h *CompletedDonationsHandler) Handle(ctx context.Context, userID string) ([]*GetCompletedDonationsResult, error) {
	donorPets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll:          true,
		IgnoreSoftDelete: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets for user")
	}

	petIDs := make([]string, len(donorPets))
	petsByID := make(map[string]*petmodel.Pet, len(donorPets))
	for i, dPet := range donorPets {
		petIDs[i] = dPet.ID
		petsByID[dPet.ID] = dPet
	}

	applicationsMap, err := h.donorRepo.GetByPetIDs(ctx, petIDs, true)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	// Собираем все неактивные отклики одним проходом, чтобы затем одним
	// batch-запросом достать связанные BloodRequest без N+1.
	var inactiveApps []*donormodel.DonorResponse
	appIDs := make([]string, 0)
	for _, applications := range applicationsMap {
		for _, app := range applications {
			if app.IsInactive() {
				inactiveApps = append(inactiveApps, app)
				appIDs = append(appIDs, app.ID)
			}
		}
	}

	bloodReqsMap, err := h.bloodReqRepo.GetByApplicationIDs(ctx, appIDs, true)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	result := make([]*GetCompletedDonationsResult, 0, len(inactiveApps))
	for _, app := range inactiveApps {
		enriched, err := h.enrichCompletedDonation(ctx, petsByID[app.DonorID], app, bloodReqsMap[app.ID])
		if err != nil {
			if errors.Is(err, apperrors.ErrBloodRequestNotFound) {
				continue
			}
			return nil, err
		}
		result = append(result, enriched)
	}

	return result, nil
}

// enrichCompletedDonation достраивает результат для одного неактивного отклика:
// питомца-реципиента и бонусы. request может быть nil, если для этого отклика
// не нашлось BloodRequest в batch-выборке — такой отклик пропускается.
func (h *CompletedDonationsHandler) enrichCompletedDonation(
	ctx context.Context,
	donorPet *petmodel.Pet,
	app *donormodel.DonorResponse,
	request *bloodreqmodel.BloodRequestWithApplications,
) (*GetCompletedDonationsResult, error) {
	if request == nil {
		return nil, apperrors.ErrBloodRequestNotFound
	}

	recipientPet, err := h.petRepo.GetByID(ctx, request.BloodRequest.PetID, pet.PetPreloadOptions{IgnoreSoftDelete: true})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pet")
	}

	bonuses, err := h.bonusRepo.GetBonusesByDonorResponseID(ctx, app.ID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get bonuses")
	}
	if len(bonuses) == 0 {
		bonuses = []*bonusmodel.Bonus{bonusmodel.NewLockBonus()}
	}

	return &GetCompletedDonationsResult{
		ApplicationData:  *app,
		BloodSearchData:  *request,
		RecipientPetData: *recipientPet,
		DonorPetData:     *donorPet,
		Bonuses:          bonuses,
	}, nil
}
