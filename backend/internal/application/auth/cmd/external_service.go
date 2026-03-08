package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	auth "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ExternalAuthHandler struct {
	userRepo       user.Repository
	appValidator   auth.AppValidator
	tokenGenerator auth.TokenGenerator
	txManager      *presistance.TxManager
}

func NewExternalSignInHandler(userepo user.Repository, appValidator auth.AppValidator, tokenGenerator auth.TokenGenerator, txManager *presistance.TxManager) *ExternalAuthHandler {
	return &ExternalAuthHandler{
		userRepo:       userepo,
		appValidator:   appValidator,
		tokenGenerator: tokenGenerator,
		txManager:      txManager,
	}
}

func (h *ExternalAuthHandler) Handle(ctx context.Context, authreq *authmodel.Identity, userdata *usermodel.User, metadata *authmodel.Metadata) (*authmodel.Identity, error) {
	authreq.ProviderUserID, authreq.ProviderName = h.appValidator.ValidateBySecret(authreq.ProviderUserID, authreq.ServiceKey)
	if authreq.ProviderUserID == 0 {
		return nil, apperrors.ErrInvalidUserData
	}
	exist, err := h.userRepo.ExistsByProvider(ctx, authreq.ProviderUserID, string(authreq.ProviderName))
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
		authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(authreq.ProviderName))
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

	authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(authreq.ProviderName))
	if err != nil {
		return nil, err
	}

	authData.AccessToken = h.tokenGenerator.Generate(authData.UserID, "user")

	return authData, nil
}
