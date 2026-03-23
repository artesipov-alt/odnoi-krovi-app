package cmd

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type CreateRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	publisher ports.EventPublisher
}

func NewCreateRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	publisher ports.EventPublisher,
) *CreateRequestHandler {
	return &CreateRequestHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		publisher: publisher,
	}
}

func (h *CreateRequestHandler) Handle(ctx context.Context, req *model.BloodRequest) (*model.BloodRequest, error) {
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

	if err := h.publisher.PublishBloodRequestCreated(ctx, events.BloodRequestCreated{
		RequestID:  newReq.ID,
		BloodTypes: req.BloodGroupNames,
		Regions:    req.Regions,
		CreatedAt:  *newReq.CreatedAt,
	}); err != nil {
		slog.Error("failed to publish blood request created event", "err", err)
		return nil, nil
	}

	return newReq, nil
}
