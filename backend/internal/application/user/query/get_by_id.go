package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type GetByIDHandler struct {
	userRepo user.Repository
}

func NewGetByIDHandler(userepo user.Repository) *GetByIDHandler {
	return &GetByIDHandler{
		userRepo: userepo,
	}
}

func (h *GetByIDHandler) Handle(ctx context.Context, id string, withPets bool, withDonorPrefs bool) (*usermodel.User, error) {
	opts := user.UserPreloadOptions{
		WithPets:            withPets,
		WithDonorPreference: withDonorPrefs,
	}

	u, err := h.userRepo.GetByID(ctx, id, opts)
	if err != nil {
		return nil, err // Доменная ошибка (например, ErrUserNotFound)
	}

	return u, nil
}
