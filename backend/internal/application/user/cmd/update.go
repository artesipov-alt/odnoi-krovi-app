package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

type UpdateHandler struct {
	userRepo user.Repository
}

func NewUpdateHandler(userRepo user.Repository) *UpdateHandler {
	return &UpdateHandler{
		userRepo: userRepo,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, id string, input *usermodel.User) (*time.Time, error) {
	err := h.userRepo.Update(ctx, id, input)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to update user")
	}

	// Return current time as UpdatedAt since repository doesn't return it
	now := time.Now()
	return &now, nil
}
