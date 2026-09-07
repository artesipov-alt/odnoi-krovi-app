package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type DeleteHandler struct {
	userRepo  user.Repository
	txManager *presistance.TxManager
}

func NewDeleteHandler(userRepo user.Repository, txManager *presistance.TxManager) *DeleteHandler {
	return &DeleteHandler{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, userID string) error {
	err := h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.userRepo.Delete(txCtx, userID); err != nil {
			return apperrors.Internal(err, "failed to delete user")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
