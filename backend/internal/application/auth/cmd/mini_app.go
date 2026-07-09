package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	auth "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type MiniAppAuthHandler struct {
	userRepo         user.Repository
	miniAppValidator auth.MiniAppValidator
	tokenGenerator   auth.TokenGenerator
	txManager        *presistance.TxManager
}

func NewMiniAppSignInHandler(userepo user.Repository, miniAppValidator auth.MiniAppValidator, tokenGenerator auth.TokenGenerator, txManager *presistance.TxManager) *MiniAppAuthHandler {
	return &MiniAppAuthHandler{
		userRepo:         userepo,
		miniAppValidator: miniAppValidator,
		tokenGenerator:   tokenGenerator,
		txManager:        txManager,
	}
}

func (h *MiniAppAuthHandler) Handle(ctx context.Context, idndata *authmodel.Identity, metadata *authmodel.Metadata) (*authmodel.Identity, error) {
	webAppData, err := h.miniAppValidator.ValidateWebAppInitData(ctx, idndata.AppInitData, idndata.ProviderName)
	if err != nil {
		return nil, apperrors.ErrInvalidUserData
	}

	if webAppData.User == nil {
		return nil, apperrors.ErrInvalidUserData
	}

	idndata.SetProviderID(webAppData.User.ID)

	existData, err := h.userRepo.GetByProvider(ctx, idndata.ProviderUserID, idndata.ProviderName)
	if err != nil {
		return nil, err
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.userRepo.UpsertUserIdentity(txCtx, existData.UserID, idndata, metadata); err != nil {
			return err
		}

		if err := h.userRepo.UpdateLastSeen(txCtx, existData.UserID); err != nil {
			return err
		}

		if metadata != nil {
			if err := h.userRepo.UpsertUTM(txCtx, existData.UserID, metadata); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	newAuthData, err := h.userRepo.GetByProvider(ctx, idndata.ProviderUserID, idndata.ProviderName)
	if err != nil {
		return nil, err
	}

	accessToken, expiresAt := h.tokenGenerator.Generate(newAuthData.UserID, string(usermodel.RoleUser), time.Now())
	newAuthData.SetJWTData(accessToken, expiresAt)

	return newAuthData, nil
}
