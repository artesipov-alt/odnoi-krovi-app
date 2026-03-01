package query

import (
	"context"
	"strconv"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type GetByIDHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	storage   filestorage.Repository
}

func NewGetByIDHandler(bloodRepo bloodsearch.BloodRequestRepository, storage filestorage.Repository) *GetByIDHandler {
	return &GetByIDHandler{
		bloodRepo: bloodRepo,
		storage:   storage,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, id string) (*ent.BloodSearchRequest, error) {
	req, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoUrls = h.buildFullPhotoURLs(req.PhotoUrls, req.UpdatedAt)

	return req, nil
}

// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
func (h *GetByIDHandler) buildFullPhotoURLs(paths []string, updatedAt time.Time) []string {
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
