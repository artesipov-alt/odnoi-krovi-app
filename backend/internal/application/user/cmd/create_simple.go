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

func (h *CreateSimpleHandler) Handle(ctx context.Context, user *usermodel.User, prefs *usermodel.DonorPreference, metadata *usermodel.Metadata) (*usermodel.User, error) {
	// Проверка exists — это координация, не бизнес-логика
	exists, _ := h.userRepo.ExistsProvider(ctx, user.ProviderID, user.ProviderName)
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	newuser, err := h.userRepo.Create(ctx, user, prefs)
	if err != nil {
		return nil, err
	}

	if err := h.userRepo.SaveUTM(ctx, newuser.ID,
		&metadata.UTMData.Source, &metadata.UTMData.Medium,
		&metadata.UTMData.Campaign,
		&metadata.UTMData.Content,
		&metadata.UTMData.Term); err != nil {
		return nil, err
	}

	return newuser, err
}
