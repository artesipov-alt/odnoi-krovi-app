package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type ApplyForRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	donorRepo bloodsearch.DonorResponseRepository
}

func NewApplyForRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	donorRepo bloodsearch.DonorResponseRepository,
) *ApplyForRequestHandler {
	return &ApplyForRequestHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		donorRepo: donorRepo,
	}
}

func (h *ApplyForRequestHandler) Handle(ctx context.Context, reqID, donorID string, conditions []string) (*model.DonorResponse, error) {
	// Проверяем существование и статус заявки
	req, err := h.bloodRepo.GetByID(ctx, reqID)
	if err != nil {
		return nil, err
	}
	if req.Status != model.BloodRequestStatusActive {
		return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
	}

	// Проверяем существование донора
	donorExists, err := h.petRepo.ExistsByID(ctx, donorID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check donor existence")
	}
	if !donorExists {
		return nil, apperrors.ErrPetNotFound
	}

	// Проверяем, нет ли уже отклика от этого донора на эту заявку
	responses, err := h.donorRepo.GetDonorResponsesByDonorID(ctx, donorID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check existing responses")
	}
	for _, resp := range responses {
		if resp.RequestID == reqID {
			return nil, apperrors.ErrDonorResponseAlreadyExists
		}
	}

	// Create domain model using constructor
	resp, err := model.NewDonorResponse(reqID, donorID, conditions)
	if err != nil {
		return nil, apperrors.Validation(err.Error(), map[string]any{"field": "donor_response"})
	}

	created, err := h.donorRepo.CreateDonorResponse(ctx, resp)
	if err != nil {
		return nil, err
	}

	return created, nil
}
