package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CreateSimpleHandler struct {
	userRepo  user.Repository
	txManager *presistance.TxManager
}

func NewCreateSimpleHandler(userepo user.Repository, txManager *presistance.TxManager) *CreateSimpleHandler {
	return &CreateSimpleHandler{
		userRepo:  userepo,
		txManager: txManager,
	}
}

func (h *CreateSimpleHandler) Handle(ctx context.Context, user *usermodel.User, prefs *usermodel.DonorPreference, metadata *usermodel.Metadata) (*usermodel.User, error) {
	// Проверка exists — это координация, не бизнес-логика
	exists, _ := h.userRepo.ExistsProvider(ctx, user.ProviderID, user.ProviderName)
	if exists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	var newuser *usermodel.User
	err := h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		var err error

		// 1. Создаем пользователя с identity
		newuser, err = h.userRepo.CreateUserWithIdentity(txCtx, user)
		if err != nil {
			return err
		}

		// 2. Создаем настройки донора
		if prefs != nil {
			if err := h.userRepo.CreateDonorPreference(txCtx, newuser.ID, prefs); err != nil {
				return err
			}
		}

		// 3. Сохраняем UTM-метки
		if metadata != nil && metadata.UTMData != nil {
			if err := h.userRepo.SaveUTM(txCtx, newuser.ID,
				&metadata.UTMData.Source, &metadata.UTMData.Medium,
				&metadata.UTMData.Campaign,
				&metadata.UTMData.Content,
				&metadata.UTMData.Term); err != nil {
				return err

			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newuser, nil
}
