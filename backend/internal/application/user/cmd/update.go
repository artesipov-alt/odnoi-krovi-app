package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/user/model"
)

type UpdateHandler struct {
	userRepo user.Repository
}

func NewUpdateHandler(userRepo user.Repository) *UpdateHandler {
	return &UpdateHandler{
		userRepo: userRepo,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, id string, input *usermodel.User) error {
	err := h.userRepo.Update(ctx, id, input)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.Internal(err, "failed to update user")
	}
	return nil
}
