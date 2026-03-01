package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
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

func (h *ApplyForRequestHandler) Handle(ctx context.Context, reqID, donorID string, conditions []string) (*ent.DonorResponse, error) {
	// Проверяем существование и статус заявки
	req, err := h.bloodRepo.GetByID(ctx, reqID)
	if err != nil {
		return nil, err
	}
	if req.Status != bloodsearchrequest.StatusActive {
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
		if resp.Edges.Request.ID == reqID {
			return nil, apperrors.ErrDonorResponseAlreadyExists
		}
	}

	resp, err := h.donorRepo.CreateDonorResponse(ctx, reqID, donorID, conditions)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
