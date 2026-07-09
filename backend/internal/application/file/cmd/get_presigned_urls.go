package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type GetPresignedURLsHandler struct {
	storage   filestorage.Repository
	petRepo   pet.Repository
	userRepo  user.Repository
	bloodRepo bloodsearch.Repository
}

func NewGetPresignedURLsHandler(
	storage filestorage.Repository,
	petRepo pet.Repository,
	userRepo user.Repository,
	bloodRepo bloodsearch.Repository,
) *GetPresignedURLsHandler {
	return &GetPresignedURLsHandler{
		storage:   storage,
		petRepo:   petRepo,
		userRepo:  userRepo,
		bloodRepo: bloodRepo,
	}
}

func (h *GetPresignedURLsHandler) Handle(ctx context.Context, id string, count int64, preloadType string) ([]filestorage.UploadInfo, error) {
	var uploadInfos []filestorage.UploadInfo

	switch preloadType {
	case "pet_avatar":
		exists, err := h.petRepo.ExistsByID(ctx, id)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check pet existence")
		}
		if !exists {
			return nil, apperrors.ErrPetNotFound
		}
		uploadInfos, err = h.storage.GetPresignedURLs(ctx, count, id)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to get upload URLs for pet")
		}
	case "user_avatar":
		exists, err := h.userRepo.ExistsByID(ctx, id)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check user existence")
		}
		if !exists {
			return nil, apperrors.ErrUserNotFound
		}
		uploadInfos, err = h.storage.GetPresignedURLs(ctx, count, id)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to get upload URLs for user")
		}
	case "blood_req":
		exists, err := h.bloodRepo.ExistsByID(ctx, id)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check blood request existence")
		}
		if !exists {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		uploadInfos, err = h.storage.GetPresignedURLs(ctx, count, id)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to get upload URLs for blood request")
		}
	default:
		return nil, apperrors.BadRequest("unsupported preload type")
	}

	if len(uploadInfos) == 0 {
		return nil, apperrors.Internal(nil, "no upload info returned")
	}

	return uploadInfos, nil
}
