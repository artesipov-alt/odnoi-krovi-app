package query

import (
	"context"

	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
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

func (h *GetByTelegramHandler) Handle(ctx context.Context, id int64, withPets bool) (*usermodel.User, []*petmodel.Pet, error) {
	opts := user.UserPreloadOptions{
		WithPets: withPets,
	}

	u, p, err := h.userRepo.GetByTelegram(ctx, id, opts)
	if err != nil {
		return nil, nil, err // Доменная ошибка (например, ErrUserNotFound)
	}

	return u, p, nil
}
