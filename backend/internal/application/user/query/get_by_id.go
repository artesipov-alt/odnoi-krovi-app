package query

import (
	"context"

	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
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

func (h *GetByIDHandler) Handle(ctx context.Context, id string, withPets bool) (*usermodel.User, []*petmodel.Pet, error) {
	opts := user.UserPreloadOptions{
		WithPets: withPets,
	}

	u, p, err := h.userRepo.GetByID(ctx, id, opts)
	if err != nil {
		return nil, nil, err // Доменная ошибка (например, ErrUserNotFound)
	}

	return u, p, nil
}
