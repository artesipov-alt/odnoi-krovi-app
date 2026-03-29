package query

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	recipientmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/recipient/model"
)

type ListRequestsHandler struct {
	bloodRepo     bloodsearch.BloodRequestRepository
	petRepo       pet.Repository
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.BloodRequestRepository
}

func NewListRequestsHandler(bloodRepo bloodsearch.BloodRequestRepository, petRepo pet.Repository) *ListRequestsHandler {
	return &ListRequestsHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*recipientmodel.Recipient, error) {
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

	var requestsWithDonors []*recipientmodel.Recipient
	for _, request := range allRequests {
		if len(request.MatchingDonors) != 0 {
			requestsWithDonors = append(requestsWithDonors, request)
		}
	}

	return requestsWithDonors, nil
}
