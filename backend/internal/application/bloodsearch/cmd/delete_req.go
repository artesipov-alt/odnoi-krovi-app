package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type DeleteRequestHandler struct {
	bloodRepo bloodsearch.Repository
	txManager *presistance.TxManager
}

func NewDeleteRequestHandler(bloodRepo bloodsearch.Repository, txManager *presistance.TxManager) *DeleteRequestHandler {
	return &DeleteRequestHandler{
		bloodRepo: bloodRepo,
		txManager: txManager,
	}
}

func (h *DeleteRequestHandler) Handle(ctx context.Context, id string) error {
	// Получаем заявку, чтобы узнать PetID
	_, err := h.bloodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		entTx := ent.TxFromContext(txCtx)
		if entTx == nil {
			return apperrors.Internal(nil, "ent.TxFromContext returned nil")
		}

		// Удаляем заявку
		err = h.bloodRepo.Delete(txCtx, id)
		if err != nil {
			if ent.IsNotFound(err) {
				return apperrors.ErrBloodRequestNotFound
			}
			return apperrors.Internal(err, "failed to delete blood request")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
