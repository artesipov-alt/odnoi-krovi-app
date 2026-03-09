package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	auth "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type MiniAppAuthHandler struct {
	userRepo       user.Repository
	maxValidator   auth.AppValidator
	tgValidator    auth.AppValidator
	tokenGenerator auth.TokenGenerator
	txManager      *presistance.TxManager
}

func NewMiniAppSignInHandler(userepo user.Repository, tgValidator auth.AppValidator, maxValidator auth.AppValidator, tokenGenerator auth.TokenGenerator, txManager *presistance.TxManager) *MiniAppAuthHandler {
	return &MiniAppAuthHandler{
		userRepo:       userepo,
		maxValidator:   maxValidator,
		tgValidator:    tgValidator,
		tokenGenerator: tokenGenerator,
		txManager:      txManager,
	}
}

func (h *MiniAppAuthHandler) Handle(ctx context.Context, authreq *authmodel.Identity, metadata *authmodel.Metadata) (*authmodel.Identity, error) {
	var validator auth.AppValidator
	providerName := authreq.ProviderName

	switch providerName {
	case authmodel.ProviderMax:
		validator = h.maxValidator
	case authmodel.ProviderTelegram:
		validator = h.tgValidator
	default:
		return nil, apperrors.ErrInvalidUserData
	}

	// Валидируем Web App initData
	webAppData, err := validator.ValidateWebAppInitData(ctx, authreq.AppInitData)
	if err != nil {
		return nil, apperrors.ErrInvalidUserData
	}

	if webAppData.User == nil {
		return nil, apperrors.ErrInvalidUserData
	}
	authreq.ProviderUserID = webAppData.User.ID

	role := "user"

	authData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(authreq.ProviderName))
	if err != nil {
		return nil, err
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		authreq.UserID = authData.UserID
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

	newAuthData, err := h.userRepo.GetByProvider(ctx, authreq.ProviderUserID, string(authreq.ProviderName))
	if err != nil {
		return nil, err
	}

	accessToken, expiresAt := h.tokenGenerator.Generate(newAuthData.UserID, role, time.Now())
	newAuthData.AccessToken = accessToken
	newAuthData.ExpiresAt = expiresAt

	return newAuthData, nil
}
