package cmd

import (
	"context"
	"strconv"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
)

type CreateRequestHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
	storage   filestorage.Repository
}

func NewCreateRequestHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	petRepo pet.Repository,
	storage filestorage.Repository,
) *CreateRequestHandler {
	return &CreateRequestHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
		storage:   storage,
	}
}

func (h *CreateRequestHandler) Handle(ctx context.Context, input *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error) {
	// Проверяем существование питомца
	exists, err := h.petRepo.ExistsByID(ctx, input.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return nil, apperrors.ErrPetNotFound
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := h.bloodRepo.ExistsByPetID(ctx, input.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check request existence")
	}
	if activeExists {
		return nil, apperrors.ErrBloodRequestAlreadyExists
	}

	status := bloodsearchrequest.StatusActive
	input.Status = &status

	newReq, err := h.bloodRepo.Create(ctx, input)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create blood request")
	}

	newReq.PhotoUrls = h.buildFullPhotoURLs(newReq.PhotoUrls, newReq.UpdatedAt)

	return newReq, nil
}

func (h *CreateRequestHandler) buildFullPhotoURLs(paths []string, updatedAt time.Time) []string {
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
