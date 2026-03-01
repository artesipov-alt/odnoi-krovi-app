package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type DeleteHandler struct {
	userRepo user.Repository
}

func NewDeleteHandler(userRepo user.Repository) *DeleteHandler {
	return &DeleteHandler{
		userRepo: userRepo,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, userID string) error {
	if err := h.userRepo.Delete(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to delete user")
	}
	return nil
}
