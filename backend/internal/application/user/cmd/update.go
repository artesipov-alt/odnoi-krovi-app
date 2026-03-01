package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type UpdateHandler struct {
	userRepo user.Repository
	storage  filestorage.Repository
}

func NewUpdateHandler(userRepo user.Repository, storage filestorage.Repository) *UpdateHandler {
	return &UpdateHandler{
		userRepo: userRepo,
		storage:  storage,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, id string, input *usermodel.User) error {
	// Нормализуем PhotoURLs - преобразуем полные URL обратно в относительные пути
	for i, url := range input.PhotoURLs {
		input.PhotoURLs[i] = h.storage.ExtractPathFromURL(url)
	}

	err := h.userRepo.Update(ctx, id, input)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.Internal(err, "failed to update user")
	}
	return nil
}
