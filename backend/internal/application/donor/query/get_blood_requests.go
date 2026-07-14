package query

import (
	"context"
	"slices"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type ListRequestsHandler struct {
	petRepo       pet.Repository
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.Repository
	matchingSvc   bloodsearch.MatchingService
	userRepo      user.Repository
}

func NewListRequestsHandler(petRepo pet.Repository, donorRespRepo donor.Repository, bloodReqRepo bloodsearch.Repository, matchingSvc bloodsearch.MatchingService, userRepo user.Repository) *ListRequestsHandler {
	return &ListRequestsHandler{
		petRepo:       petRepo,
		donorRespRepo: donorRespRepo,
		bloodReqRepo:  bloodReqRepo,
		matchingSvc:   matchingSvc,
		userRepo:      userRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithMatchingDonors, error) {
	user, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{
		WithDonorPreference: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get user")
	}

	if user.DonorPreference != nil {
		return nil, apperrors.BadRequest("Настройки донора не заполнены")
	}

	preferredLocations := user.DonorPreference.PreferredLocationIDs

	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	// Collect pet IDs for batch queries
	petIDs := make([]string, len(pets))
	for i, pet := range pets {
		petIDs[i] = pet.ID
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := h.donorRespRepo.GetByPetIDs(ctx, petIDs, false)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodReqRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	for _, recipientPet := range pets {
		applications := applicationsMap[recipientPet.ID]
		var application *donormodel.DonorResponse
		for _, app := range applications {
			if app.IsActiveForDonation() {
				application = app
				break
			}
		}
		bloodReq := bloodReqsMap[recipientPet.ID]
		recipientPet.RecalculateStatus(time.Now(), pet.BuildDonationContext(application, bloodReq))
	}

	potentialDonors := petmodel.FilterDonors(pets)

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
