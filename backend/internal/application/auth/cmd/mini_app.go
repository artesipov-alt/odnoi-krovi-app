package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	auth "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type MiniAppAuthHandler struct {
	userRepo       user.Repository
	appValidator   auth.AppValidator
	tokenGenerator auth.TokenGenerator
	txManager      *presistance.TxManager
}

func NewMiniAppSignInHandler(userepo user.Repository, appValidator auth.AppValidator, tokenGenerator auth.TokenGenerator, txManager *presistance.TxManager) *MiniAppAuthHandler {
	return &MiniAppAuthHandler{
		userRepo:       userepo,
		appValidator:   appValidator,
		tokenGenerator: tokenGenerator,
		txManager:      txManager,
	}
}

func (h *MiniAppAuthHandler) Handle(ctx context.Context, authreq *authmodel.Identity, metadata *authmodel.Metadata) (*authmodel.Identity, error) {
	var role string
	authreq.ProviderUserID, role = h.appValidator.ValidateHash(authreq.AppInitData)
	if authreq.ProviderUserID == "" {
		return nil, apperrors.ErrInvalidUserData
	}

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

	newAuthData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(authreq.ProviderName))
	if err != nil {
		return nil, err
	}

	authData.AccessToken = h.tokenGenerator.Generate(newAuthData.UserID, role)

	return authData, nil
}
