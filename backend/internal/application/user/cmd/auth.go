package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type AuthHandler struct {
	userRepo  user.Repository
	txManager *presistance.TxManager
}

func NewAuthHandler(userepo user.Repository, txManager *presistance.TxManager) *AuthHandler {
	return &AuthHandler{
		userRepo:  userepo,
		txManager: txManager,
	}
}

func (h *AuthHandler) Handle(ctx context.Context, authreq *usermodel.Identity, userdata *usermodel.User, metadata *usermodel.Metadata) (*usermodel.Identity, error) {
	exist, err := h.userRepo.ExistsByProvider(ctx, authreq.ProviderUserID, authreq.ProviderName)
	if err != nil {
		return nil, err
	}
	if !exist {
		err := h.txManager.WithTx(ctx, func(txCtx context.Context) error {
			// 1. Создаем пользователя
			newuser, err := h.userRepo.CreateUser(txCtx, userdata)
			if err != nil {
				return err
			}
			authreq.UserID = newuser.ID
			// 2. Создаем identity пользователя
			if err := h.userRepo.UpsertUserIdentity(txCtx, authreq); err != nil {
				return err
			}

			if metadata != nil {
				if err := h.userRepo.UpsertUTM(txCtx, newuser.ID, metadata); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Пользователь существует - обновляем метаданные и UTM
		authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, authreq.ProviderName)
		if err != nil {
			return nil, err
		}

		err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
			// Обновляем identity с новыми метаданными
			authreq.UserID = authData.UserID
			if err := h.userRepo.UpsertUserIdentity(txCtx, authreq); err != nil {
				return err
			}

			// Обновляем/создаем UTM метки
			if metadata != nil {
				if err := h.userRepo.UpsertUTM(txCtx, authData.UserID, metadata); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, authreq.ProviderName)
	if err != nil {
		return nil, err
	}

	authData.Token = "secret-token-" + authreq.AppInitData

	return authData, nil
}
