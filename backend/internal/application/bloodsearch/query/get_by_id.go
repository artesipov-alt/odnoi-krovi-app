package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
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

func (h *GetByIDHandler) Handle(ctx context.Context, id string) (*model.BloodRequest, error) {
	req, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Преобразуем пути к фото в полные URL
	req.PhotoURLs = h.storage.BuildPhotoURLs(req.PhotoURLs, req.UpdatedAt)

	return req, nil
}
