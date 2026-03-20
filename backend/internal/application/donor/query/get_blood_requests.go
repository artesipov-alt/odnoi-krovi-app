package query

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type ListRequestsHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
}

func NewListRequestsHandler(bloodRepo bloodsearch.BloodRequestRepository, petRepo pet.Repository) *ListRequestsHandler {
	return &ListRequestsHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*donormodel.Recipient, error) {
	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll:      true,
		WithBloodReq: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}
	now := time.Now()
	for i, _ := range pets {
		pets[i].RecalculateFactors(now)
		pets[i].CalculateDonorStatus()
	}

	potentialDonors := petmodel.FilterDonors(pets)

	allRequests, err := h.bloodRepo.AdaptiveList(ctx, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	// Find matching donors
	for _, recipient := range allRequests {
		for _, donor := range potentialDonors {
			recipient.MatchDonor(donor)
		}
	}

	var requestsWithDonors []*donormodel.Recipient
	for _, request := range allRequests {
		if len(request.MatchingDonors) != 0 {
			requestsWithDonors = append(requestsWithDonors, request)
		}
	}

	return requestsWithDonors, nil
}
