package cmd

import (
	"context"
	"strconv"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
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

func (h *UpdateRequestHandler) Handle(ctx context.Context, id string, bloodReq *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
	// Проверяем существование
	req, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем текущий статус, если он не передан
	if bloodReq.Status == nil {
		bloodReq.Status = &req.Status
	}

	updatedReq, err := h.bloodRepo.Update(ctx, id, bloodReq)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	updatedReq.PhotoUrls = h.buildFullPhotoURLs(updatedReq.PhotoUrls, updatedReq.UpdatedAt)

	return updatedReq, nil
}

// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (h *UpdateRequestHandler) buildFullPhotoURLs(paths []string, updatedAt time.Time) []string {
	if len(paths) == 0 {
		return []string{}
	}
	result := make([]string, len(paths))
	for i, path := range paths {
		if path == "" {
			result[i] = ""
		} else {
			url := h.storage.GetPublicURLFromPath(path)
			result[i] = url + "?t=" + strconv.FormatInt(updatedAt.Unix(), 10)
		}
	}
	return result
}
