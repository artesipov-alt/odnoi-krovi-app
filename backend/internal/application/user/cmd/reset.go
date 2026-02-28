package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user"
)

type ResetHandler struct {
	userRepo user.Repository
}

func NewResetHandler(userRepo user.Repository) *ResetHandler {
	return &ResetHandler{
		userRepo: userRepo,
	}
}

func (h *ResetHandler) Handle(ctx context.Context, userID string) error {
	if err := h.userRepo.ResetUser(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to reset user")
	}
	return nil
}
