package query

import (
	"context"
	"slices"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type ListRequestsHandler struct {
	petRepo      pet.Repository
	bloodReqRepo bloodsearch.Repository
	matchingSvc  bloodsearch.MatchingService
	enricher     enrich.PetEnricher
	userRepo     user.Repository
}

func NewListRequestsHandler(petRepo pet.Repository, bloodReqRepo bloodsearch.Repository, matchingSvc bloodsearch.MatchingService, enricher enrich.PetEnricher, userRepo user.Repository) *ListRequestsHandler {
	return &ListRequestsHandler{
		petRepo:      petRepo,
		bloodReqRepo: bloodReqRepo,
		matchingSvc:  matchingSvc,
		enricher:     enricher,
		userRepo:     userRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithMatchingDonors, error) {
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

	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	if _, err := h.enricher.RecalculateAll(ctx, pets, enrich.Options{RecoveryPeriodMonths: user.DonorPreference.RecoveryPeriodMonths}); err != nil {
		return nil, apperrors.Internal(err, "failed to recalculate all pets")
	}

	potentialDonors := petmodel.FilterDonors(pets)
	if len(potentialDonors) == 0 {
		return []*bloodreqmodel.BloodRequestWithMatchingDonors{}, nil
	}

	allRequests, err := h.bloodReqRepo.AdaptiveList(ctx, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	// Find matching donors
	for _, recipient := range allRequests {
		recipient.SyncPrivilegeAndPriority()
		for _, donor := range potentialDonors {
			if d, ok := h.matchingSvc.MatchDonor(recipient, donor, preferredLocations); ok {
				recipient.AddMatchingDonor(d)
			}
		}
	}

	var requestsWithDonors []*bloodreqmodel.BloodRequestWithMatchingDonors
	for _, request := range allRequests {
		if len(request.MatchingDonors) != 0 {
			requestsWithDonors = append(requestsWithDonors, request)
		}
	}

	slices.SortFunc(requestsWithDonors, sortByPriorityAndDate)

	return requestsWithDonors, nil
}

// Сортируем массив реципиентов по дате создания от старых к новым и по приоритету
func sortByPriorityAndDate(a, b *bloodreqmodel.BloodRequestWithMatchingDonors) int {
	if a.BloodRequest.PrioritySearch == b.BloodRequest.PrioritySearch {
		return b.BloodRequest.CreatedAt.Compare(*a.BloodRequest.CreatedAt)
	}
	if a.BloodRequest.PrioritySearch {
		return -1
	}
	return 1
}
