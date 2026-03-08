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
	// Валидируем и получаем provider name
	_, providerName, role := h.appValidator.ValidateBySecret(ctx, authreq.ProviderUserID, authreq.ServiceKey)
	if providerName == "" {
		return nil, apperrors.ErrInvalidUserData
	}
	authreq.ProviderName = authmodel.ProviderName(providerName)

	exist, err := h.userRepo.ExistsByProvider(ctx, authreq.ProviderUserID, providerName)
	if err != nil {
		return nil, err
	}

	if !exist {
		// Создаем нового пользователя
		err := h.txManager.WithTx(ctx, func(txCtx context.Context) error {
			newuser, err := h.userRepo.CreateUser(txCtx, userdata)
			if err != nil {
				return err
			}
			authreq.UserID = newuser.ID
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
		authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(providerName))
		if err != nil {
			return nil, err
		}

		err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
			authreq.UserID = authData.UserID
			// Создаем identity пользователя
			if err := h.userRepo.UpsertUserIdentity(txCtx, authreq); err != nil {
				return err
			}
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

	authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(providerName))
	if err != nil {
		return nil, err
	}

	authData.AccessToken = h.tokenGenerator.Generate(authData.UserID, role)
	return authData, nil
}
