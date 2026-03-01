package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
)

type GetByPetIDHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	storage   filestorage.Repository
}

func NewGetByPetIDHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	storage filestorage.Repository,
) *GetByPetIDHandler {
	return &GetByPetIDHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		storage:   storage,
	}
}

func (h *GetByPetIDHandler) Handle(ctx context.Context, petID string) (*model.BloodRequest, int, error) {
	req, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, 0, err
	}

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, req.BloodGroupNames)
	if err != nil {
		return nil, 0, err
	}

	req.PhotoURLs = h.storage.BuildPhotoURLs(req.PhotoURLs, req.UpdatedAt)

	return req, suitableDonors, nil
}
