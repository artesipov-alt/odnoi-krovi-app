package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type GetByTelegramHandler struct {
	userRepo user.Repository
}

func NewGetByTelegramHandler(userepo user.Repository) *GetByTelegramHandler {
	return &GetByTelegramHandler{
		userRepo: userepo,
	}
}

func (h *GetByTelegramHandler) Handle(ctx context.Context, id int64, withPets bool, withDonorPrefs bool) (*usermodel.User, error) {
	opts := user.UserPreloadOptions{
		WithPets:            withPets,
		WithDonorPreference: withDonorPrefs,
	}

	u, err := h.userRepo.GetByTelegram(ctx, id, opts)
	if err != nil {
		return nil, err // Доменная ошибка (например, ErrUserNotFound)
	}

	return u, nil
}
