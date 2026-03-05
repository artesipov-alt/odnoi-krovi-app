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

func (h *AuthHandler) Handle(ctx context.Context, authreq *usermodel.Identity) (*usermodel.Identity, error) {
	authdata, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, authreq.ProviderName)
	if err != nil {
		return nil, err
	}
	if authreq.XBToken != "" {
		authdata.XBToken = authreq.XBToken + "-encrypted_from_backend"
	}

	return authdata, nil
}
