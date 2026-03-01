package query

import (
	"context"
	"strconv"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
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

func (h *GetByPetIDHandler) Handle(ctx context.Context, petID string) (*ent.BloodSearchRequest, int, error) {
	req, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, 0, err
	}

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, req.BloodGroupNames)
	if err != nil {
		return nil, 0, err
	}

	req.PhotoUrls = h.buildFullPhotoURLs(req.PhotoUrls, req.UpdatedAt)

	return req, suitableDonors, nil
}

func (h *GetByPetIDHandler) buildFullPhotoURLs(paths []string, updatedAt time.Time) []string {
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
