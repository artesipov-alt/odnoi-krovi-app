package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user"
)

type RestoreHandler struct {
	userRepo user.Repository
}

func NewRestoreHandler(userRepo user.Repository) *RestoreHandler {
	return &RestoreHandler{
		userRepo: userRepo,
	}
}

func (h *RestoreHandler) Handle(ctx context.Context, userID string) error {
	if err := h.userRepo.RestoreUser(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to restore user")
	}
	return nil
}
