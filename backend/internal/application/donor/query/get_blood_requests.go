package query

import (
	"context"
	"errors"
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

	for _, pet := range pets {
		application, err := h.donorRespRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
			return nil, apperrors.Internal(err, "failed to get donor application")
		}
		bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			return nil, apperrors.Internal(err, "failed to get blood request")
		}
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
