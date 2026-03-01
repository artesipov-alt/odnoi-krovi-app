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

func (h *CreateSimpleHandler) Handle(ctx context.Context, telegramID int64, fullName string, role string) (*usermodel.User, error) {
	// Проверка exists — это координация, не бизнес-логика
	exists, _ := h.userRepo.ExistsByTelegramID(ctx, telegramID)
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	if fullName == "" {
		fullName = "Пользователь Telegram"
	}

	if role == "" {
		role = "user"
	}

	user := &usermodel.User{
		TelegramID: telegramID,
		FullName:   fullName,
		Role:       role,
	}

	return h.userRepo.Create(ctx, user)
}
