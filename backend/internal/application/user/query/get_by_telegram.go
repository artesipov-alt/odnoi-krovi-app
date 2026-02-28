package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/filestorage"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/user/model"
)

type GetByTelegramHandler struct {
	userRepo    user.Repository
	fileStorage filestorage.Repository // Интерфейс для URL
}

func NewGetByTelegramHandler(userepo user.Repository, filestorage filestorage.Repository) *GetByTelegramHandler {
	return &GetByTelegramHandler{
		userRepo:    userepo,
		fileStorage: filestorage,
	}
}

func (h *GetByTelegramHandler) Handle(ctx context.Context, id int64, withPets bool) (*usermodel.User, []*petmodel.Pet, error) {
	opts := user.UserPreloadOptions{
		WithPets: withPets,
	}

	u, p, err := h.userRepo.GetByTelegram(ctx, id, opts)
	if err != nil {
		return nil, nil, err // Доменная ошибка (например, ErrUserNotFound)
	}

	u.PhotoURLs = h.fileStorage.BuildPhotoURLs(u.PhotoURLs, *u.UpdatedAt)

	return u, p, nil
}
