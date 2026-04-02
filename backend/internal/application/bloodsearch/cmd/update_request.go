package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
)

type UpdateRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
}

func NewUpdateRequestHandler(bloodRepo bloodsearch.BloodRequestRepository) *UpdateRequestHandler {
	return &UpdateRequestHandler{
		bloodRepo: bloodRepo,
	}
}

func (h *UpdateRequestHandler) Handle(ctx context.Context, id string, bloodReq *model.BloodRequestWithApplications) (*model.BloodRequestWithApplications, error) {
	// Проверяем существование
	existingReq, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Если статус не передан, сохраняем текущий
	if bloodReq.Status == "" {
		bloodReq.Status = existingReq.Status
	}

	updatedReq, err := h.bloodRepo.Update(ctx, id, &bloodReq.BloodRequest)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	return updatedReq, nil
}
