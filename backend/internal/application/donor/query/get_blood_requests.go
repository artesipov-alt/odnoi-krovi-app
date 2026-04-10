package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type ListRequestsHandler struct {
	petRepo       pet.Repository
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.BloodRequestRepository
	matchingSvc   bloodsearch.MatchingService
}

func NewListRequestsHandler(petRepo pet.Repository, donorRespRepo donor.Repository, bloodReqRepo bloodsearch.BloodRequestRepository, matchingSvc bloodsearch.MatchingService) *ListRequestsHandler {
	return &ListRequestsHandler{
		petRepo:       petRepo,
		donorRespRepo: donorRespRepo,
		bloodReqRepo:  bloodReqRepo,
		matchingSvc:   matchingSvc,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithMatchingDonors, error) {
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
	applicationsMap, err := h.donorRespRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodReqRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	for _, pet := range pets {
		applications := applicationsMap[pet.ID]
		application := findActiveApplication(applications)
		bloodReq := bloodReqsMap[pet.ID]
		pet.RecalculateFactors(time.Now(), application, bloodReq)
		pet.CalculateStatus(application, bloodReq)
	}

	potentialDonors := petmodel.FilterDonors(pets)

	allRequests, err := h.bloodReqRepo.AdaptiveList(ctx, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	// Find matching donors
	for _, recipient := range allRequests {
		for _, donor := range potentialDonors {
			h.matchingSvc.MatchDonor(recipient, donor)
		}
	}

	var requestsWithDonors []*bloodreqmodel.BloodRequestWithMatchingDonors
	for _, request := range allRequests {
		if len(request.MatchingDonors) != 0 {
			requestsWithDonors = append(requestsWithDonors, request)
		}
	}

	return requestsWithDonors, nil
}
