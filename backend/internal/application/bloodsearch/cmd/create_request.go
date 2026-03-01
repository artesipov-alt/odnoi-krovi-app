package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
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

func (h *CreateRequestHandler) Handle(ctx context.Context, req *model.BloodRequest) (*model.BloodRequest, error) {
	// Проверяем существование питомца
	exists, err := h.petRepo.ExistsByID(ctx, req.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check pet existence")
	}
	if !exists {
		return nil, apperrors.ErrPetNotFound
	}

	// Проверяем, нет ли уже активной заявки для этого питомца
	activeExists, err := h.bloodRepo.ExistsByPetID(ctx, req.PetID)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to check request existence")
	}
	if activeExists {
		return nil, apperrors.ErrBloodRequestAlreadyExists
	}

	// Устанавливаем статус active если не указан
	if req.Status == "" {
		req.Status = model.BloodRequestStatusActive
	}

	newReq, err := h.bloodRepo.Create(ctx, req)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create blood request")
	}

	// Строим полные URL для фото
	newReq.PhotoURLs = h.storage.BuildPhotoURLs(newReq.PhotoURLs, newReq.UpdatedAt)

	return newReq, nil
}
