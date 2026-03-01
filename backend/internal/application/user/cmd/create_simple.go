package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type CreateSimpleHandler struct {
	userRepo user.Repository
}

func NewCreateSimpleHandler(userepo user.Repository) *CreateSimpleHandler {
	return &CreateSimpleHandler{
		userRepo: userepo,
	}
}

func (h *CreateSimpleHandler) Handle(ctx context.Context, user *usermodel.User) (*usermodel.User, error) {
	// Проверка exists — это координация, не бизнес-логика
	exists, _ := h.userRepo.ExistsByTelegramID(ctx, user.TelegramID)
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	return h.userRepo.Create(ctx, user)
}
