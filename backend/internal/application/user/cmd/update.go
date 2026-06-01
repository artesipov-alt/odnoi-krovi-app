package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type UpdateHandler struct {
	userRepo  user.Repository
	txManager *presistance.TxManager
}

func NewUpdateHandler(userRepo user.Repository, txManager *presistance.TxManager) *UpdateHandler {
	return &UpdateHandler{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, id string, input *usermodel.User) (*usermodel.User, error) {
	err := h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.userRepo.UpdateUserFields(txCtx, id, input); err != nil {
			if ent.IsNotFound(err) {
				return apperrors.ErrUserNotFound
			}
			return apperrors.Internal(err, "failed to update user")
		}

		if input.DonorPreference != nil {
			if err := h.userRepo.UpsertDonorPreference(txCtx, id, input.DonorPreference); err != nil {
				return apperrors.Internal(err, "failed to upsert donor preference")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	usr, err := h.userRepo.GetByID(ctx, id, user.UserPreloadOptions{})
	if err != nil {
		return nil, err
	}

	return usr, nil
}
