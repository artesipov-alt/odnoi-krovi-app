package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type ConfirmUploadHandler struct {
	petRepo   pet.Repository
	userRepo  user.Repository
	bloodRepo bloodsearch.Repository
	storage   filestorage.Repository
}

func NewConfirmUploadHandler(
	petRepo pet.Repository,
	userRepo user.Repository,
	bloodRepo bloodsearch.Repository,
	storage filestorage.Repository,
) *ConfirmUploadHandler {
	return &ConfirmUploadHandler{
		petRepo:   petRepo,
		userRepo:  userRepo,
		bloodRepo: bloodRepo,
		storage:   storage,
	}
}

func (h *ConfirmUploadHandler) Handle(ctx context.Context, id string, paths []string, preload string) error {
	var (
		exists bool
		err    error
	)

	switch preload {
	case "pet_avatar":
		exists, err = h.petRepo.ExistsByID(ctx, id)
		if err != nil {
			return apperrors.Internal(err, "failed to check pet existence")
		}
		if !exists {
			return apperrors.ErrPetNotFound
		}
		err = h.storage.ConfirmUploads(ctx, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to confirm uploads for pet")
		}
		err = h.petRepo.AddPhotoURLs(ctx, id, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to update pet photos")
		}
	case "user_avatar":
		exists, err = h.userRepo.ExistsByID(ctx, id)
		if err != nil {
			return apperrors.Internal(err, "failed to check user existence")
		}
		if !exists {
			return apperrors.ErrUserNotFound
		}
		err = h.storage.ConfirmUploads(ctx, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to confirm uploads for user")
		}
		if err := h.userRepo.AddPhotoURLs(ctx, id, paths); err != nil {
			return apperrors.Internal(err, "failed to update user photos")
		}
	case "blood_req":
		exists, err = h.bloodRepo.ExistsByID(ctx, id)
		if err != nil {
			return apperrors.Internal(err, "failed to check blood request existence")
		}
		if !exists {
			return apperrors.ErrBloodRequestNotFound
		}
		err = h.storage.ConfirmUploads(ctx, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to confirm uploads for blood request")
		}
		if err := h.bloodRepo.AddPhotoURLs(ctx, id, paths); err != nil {
			return apperrors.Internal(err, "failed to update blood request photos")
		}
	default:
		return apperrors.BadRequest("unsupported preload type")
	}

	return nil
}
