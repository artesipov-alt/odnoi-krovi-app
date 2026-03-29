package query

import (
	"context"
	"errors"
	"log"
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

func NewListRequestsHandler(bloodRepo bloodsearch.BloodRequestRepository, petRepo pet.Repository, donorRespRepo donor.Repository, bloodReqRepo bloodsearch.BloodRequestRepository) *ListRequestsHandler {
	return &ListRequestsHandler{
		bloodRepo:     bloodRepo,
		petRepo:       petRepo,
		donorRespRepo: donorRespRepo,
		bloodReqRepo:  bloodReqRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*recipientmodel.Recipient, error) {
	log.Printf("Starting ListRequestsHandler.Handle for userID: %s", userID)
	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll: true,
	})
	if err != nil {
		log.Printf("Error getting pets for userID %s: %v", userID, err)
		return nil, apperrors.Internal(err, "failed to get pets")
	}
	log.Printf("Found %d pets for userID %s", len(pets), userID)

	for _, pet := range pets {
		log.Printf("Processing pet ID: %s", pet.ID)
		application, err := h.donorRespRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
			log.Printf("Error getting donor application for pet ID %s: %v", pet.ID, err)
			return nil, apperrors.Internal(err, "failed to get donor application")
		}
		if application != nil {
			log.Printf("Found donor application for pet ID %s", pet.ID)
		} else {
			log.Printf("No donor application found for pet ID %s", pet.ID)
		}
		bloodReq, err := h.bloodReqRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			log.Printf("Error getting blood request for pet ID %s: %v", pet.ID, err)
			return nil, apperrors.Internal(err, "failed to get blood request")
		}
		if bloodReq != nil {
			log.Printf("Found blood request for pet ID %s", pet.ID)
		} else {
			log.Printf("No blood request found for pet ID %s", pet.ID)
		}
		pet.RecalculateFactors(time.Now(), application, bloodReq)
		log.Printf("Recalculated factors for pet ID %s", pet.ID)
		pet.CalculateStatus(application, bloodReq)
		log.Printf("Calculated status for pet ID %s", pet.ID)
	}

	potentialDonors := petmodel.FilterDonors(pets)
	log.Printf("Filtered %d potential donors", len(potentialDonors))

	allRequests, err := h.bloodRepo.AdaptiveList(ctx, filters)
	if err != nil {
		log.Printf("Error listing blood requests: %v", err)
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}
	log.Printf("Found %d blood requests", len(allRequests))

	// Find matching donors
	for _, recipient := range allRequests {
		log.Printf("Matching donors for recipient ID: %s", recipient.ID)
		for _, donor := range potentialDonors {
			recipient.MatchDonor(donor)
			log.Printf("Attempted to match donor %s with recipient %s. Matches found: %d", donor.ID, recipient.ID, len(recipient.MatchingDonors))
		}
	}

	var requestsWithDonors []*recipientmodel.Recipient
	for _, request := range allRequests {
		if len(request.MatchingDonors) != 0 {
			requestsWithDonors = append(requestsWithDonors, request)
			log.Printf("Recipient %s has %d matching donors.", request.ID, len(request.MatchingDonors))
		} else {
			log.Printf("Recipient %s has no matching donors.", request.ID)
		}
	}
	log.Printf("Returning %d requests with matching donors.", len(requestsWithDonors))

	return requestsWithDonors, nil
}
