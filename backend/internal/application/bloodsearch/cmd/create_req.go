package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
)

// findActiveApplication finds the most recent active donor application from the list.
// Active means status is Accepted, Pending, or Completed but not confirmed.
func findActiveApplication(applications []*donormodel.DonorResponse) *donormodel.DonorResponse {
	for _, app := range applications {
		if app.Status == donormodel.DonorResponseStatusAccepted ||
			app.Status == donormodel.DonorResponseStatusPending ||
			(app.Status == donormodel.DonorResponseStatusCompleted && !app.IsConfirmed) {
			return app
		}
	}
	return nil
}

type CreateRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	donorRepo donor.Repository
	publisher ports.EventPublisher
}

func NewCreateRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	donorRepo donor.Repository,
	publisher ports.EventPublisher,
) *CreateRequestHandler {
	return &CreateRequestHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		donorRepo: donorRepo,
		publisher: publisher,
	}
}

func (h *CreateRequestHandler) Handle(ctx context.Context, req *model.BloodRequest) (*model.BloodRequestWithApplications, error) {
	// Проверяем существование питомца
	exists, err := h.petRepo.ExistsByID(ctx, req.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return nil, apperrors.ErrPetNotFound
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := h.bloodRepo.ExistsByPetID(ctx, req.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check request existence")
	}
	if activeExists {
		return nil, apperrors.ErrBloodRequestAlreadyExists
	}

	newReq, err := h.bloodRepo.Create(ctx, req)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create blood request")
	}

	pets, err := h.petRepo.GetPetsByBloodGroupAndRegion(ctx, req.BloodGroupNames, req.Regions)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	// Collect pet IDs for batch queries
	petIDs := make([]string, len(pets))
	for i, pet := range pets {
		petIDs[i] = pet.ID
	}

	// Batch fetch applications and blood requests
	applicationsMap, err := h.donorRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get donor applications")
	}

	bloodReqsMap, err := h.bloodRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}

	timeNow := time.Now()
	for _, pet := range pets {
		applications := applicationsMap[pet.ID]
		donorApplication := findActiveApplication(applications)
		donorBloodReq := bloodReqsMap[pet.ID]
		pet.RecalculateFactors(timeNow, donorApplication, donorBloodReq)
		pet.CalculateStatus(donorApplication, donorBloodReq)
	}

	var avilableDonors []petmodel.Pet
	for _, pet := range pets {
		if pet.PetStatus == petmodel.PetStatusDonor {
			avilableDonors = append(avilableDonors, *pet)
		}
	}

	if err := h.publisher.PublishBloodRequestCreated(ctx, events.BloodRequestCreated{
		RequestID:  newReq.ID,
		BloodTypes: req.BloodGroupNames,
		Regions:    req.Regions,
		CreatedAt:  *newReq.CreatedAt,
	}); err != nil {
		slog.Error("failed to publish blood request created event", "err", err)
		return newReq, nil
	}

	return newReq, nil
}
