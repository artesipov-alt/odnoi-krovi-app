package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type AuthHandler struct {
	userRepo user.Repository
}

func NewAuthHandler(userepo user.Repository) *AuthHandler {
	return &AuthHandler{
		userRepo: userepo,
	}
}

func (h *AuthHandler) Handle(ctx context.Context, authdata *usermodel.Identity) (*usermodel.Identity, error) {
	// // Проверка exists — это координация, не бизнес-логика
	// exists, _ := h.userRepo.ExistsByTelegramID(ctx, user.TelegramID)
	// if exists {
	// 	return nil, apperrors.ErrUserAlreadyExists
	// }

	return h.userRepo.GetByProvider(ctx, authdata.ProviderUserID, authdata.ProviderName)
}
