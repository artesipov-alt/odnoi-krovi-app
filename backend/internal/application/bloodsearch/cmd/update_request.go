package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
)

type UpdateRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	storage   filestorage.Repository
}

func NewUpdateRequestHandler(bloodRepo bloodsearch.BloodRequestRepository, storage filestorage.Repository) *UpdateRequestHandler {
	return &UpdateRequestHandler{
		bloodRepo: bloodRepo,
		storage:   storage,
	}
}

func (h *UpdateRequestHandler) Handle(ctx context.Context, id string, bloodReq *model.BloodRequest) (*model.BloodRequest, error) {
	// Проверяем существование
	existingReq, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Если статус не передан, сохраняем текущий
	if bloodReq.Status == "" {
		bloodReq.Status = existingReq.Status
	}

	// Нормализуем PhotoURLs - преобразуем полные URL обратно в относительные пути
	for i, url := range bloodReq.PhotoURLs {
		bloodReq.PhotoURLs[i] = h.storage.ExtractPathFromURL(url)
	}

	updatedReq, err := h.bloodRepo.Update(ctx, id, bloodReq)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	return updatedReq, nil
}
