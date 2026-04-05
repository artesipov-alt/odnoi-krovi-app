package cmd

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
)

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

	timeNow := time.Now()
	for _, pet := range pets {
		donorApplication, err := h.donorRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrDonorResponseNotFound) {
			return nil, apperrors.Internal(err, "failed to get donor applications")
		}
		donorBloodReq, err := h.bloodRepo.GetByPetID(ctx, pet.ID)
		if err != nil && !errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			return nil, apperrors.Internal(err, "failed to get donor blood request")
		}
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
