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

func extractUTMFromMetadata(metadata map[string]string) (string, string, string, string, string) {
	return metadata["utm_source"], metadata["utm_medium"], metadata["utm_campaign"], metadata["utm_content"], metadata["utm_term"]
}

func (h *CreateSimpleHandler) Handle(ctx context.Context, user *usermodel.User, prefs *usermodel.DonorPreference) (*usermodel.User, error) {
	// Проверка exists — это координация, не бизнес-логика
	exists, _ := h.userRepo.ExistsProvider(ctx, user.ProviderID, user.ProviderName)
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	metadata := user.MetaData

	newuser, err := h.userRepo.Create(ctx, user, prefs)
	if err != nil {
		return nil, err
	}

	source, medium, campaign, content, term := extractUTMFromMetadata(metadata)
	if err := h.userRepo.SaveUTM(ctx, newuser.ID, &source, &medium, &campaign, &content, &term); err != nil {
		return nil, err
	}

	return newuser, err
}
