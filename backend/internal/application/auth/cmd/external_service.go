package cmd

import (
	"context"
	"time"

	auth "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/partner"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ExternalAuthHandler struct {
	userRepo       user.Repository
	partnerRepo    partner.Repository
	tokenGenerator auth.TokenGenerator
	txManager      *presistance.TxManager
}

func NewExternalSignInHandler(userepo user.Repository, partnerRepo partner.Repository, tokenGenerator auth.TokenGenerator, txManager *presistance.TxManager) *ExternalAuthHandler {
	return &ExternalAuthHandler{
		userRepo:       userepo,
		partnerRepo:    partnerRepo,
		tokenGenerator: tokenGenerator,
		txManager:      txManager,
	}
}

func (h *ExternalAuthHandler) Handle(ctx context.Context, idndata *authmodel.Identity, userdata *usermodel.User, metadata *authmodel.Metadata) (*authmodel.Identity, error) {
	partner, err := h.partnerRepo.GetByAPIKey(ctx, idndata.ServiceKey)
	if err != nil {
		return nil, err
	}

	idndata.SetPartnerID(partner.ID)
	userdata.SetRole(partner.Role)

	exist, err := h.userRepo.ExistsByProvider(ctx, idndata.ProviderUserID, idndata.ProviderName)
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
			if err := h.userRepo.UpsertUserIdentity(txCtx, newuser.ID, idndata, metadata); err != nil {
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
		idn, err := h.userRepo.GetByProvider(ctx, idndata.ProviderUserID, idndata.ProviderName)
		if err != nil {
			return nil, err
		}

		err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
			idndata.SetSystemUserID(idn.UserID)
			// Создаем identity пользователя
			if err := h.userRepo.UpsertUserIdentity(txCtx, idn.UserID, idndata, metadata); err != nil {
				return err
			}
			if metadata != nil {
				if err := h.userRepo.UpsertUTM(txCtx, idn.UserID, metadata); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	idn, err := h.userRepo.GetByProvider(ctx, idndata.ProviderUserID, idndata.ProviderName)
	if err != nil {
		return nil, err
	}

	usr, err := h.userRepo.GetByID(ctx, idn.UserID, user.UserPreloadOptions{})
	if err != nil {
		return nil, err
	}

	accessToken, expiresAt := h.tokenGenerator.Generate(idn.UserID, string(usr.Role), time.Now())
	idn.SetJWTData(accessToken, expiresAt)

	return idn, nil
}
