package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type UpdateStatusHandler struct {
	txManager presistance.TxManager
	bloodRepo bloodsearch.BloodRequestRepository
}

func NewUpdateStatusHandler(txManager presistance.TxManager, bloodRepo bloodsearch.BloodRequestRepository) *UpdateStatusHandler {
	return &UpdateStatusHandler{
		txManager: txManager,
		bloodRepo: bloodRepo,
	}
}

func (h *UpdateStatusHandler) Handle(ctx context.Context, id string, status string) error {
	if err := bloodsearchrequest.StatusValidator(bloodsearchrequest.Status(status)); err != nil {
		return apperrors.ErrInvalidBloodRequestStatus.WithInternal(err)
	}

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

		// Обновляем статус заявки
		err = h.bloodRepo.UpdateStatus(txCtx, id, status)
		if err != nil {
			return apperrors.Internal(err, "failed to update status")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
