package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type CancelDonationHandler struct {
	donorRepo donor.Repository
	bloodRepo bloodsearch.BloodRequestRepository
	txManager *presistance.TxManager
}

func NewCancelDonationHandler(
	donorRepo donor.Repository,
	bloodRepo bloodsearch.BloodRequestRepository,
	txManager *presistance.TxManager,
) *CancelDonationHandler {
	return &CancelDonationHandler{
		donorRepo: donorRepo,
		bloodRepo: bloodRepo,
		txManager: txManager,
	}
}

func (h *CancelDonationHandler) Handle(ctx context.Context, resID string) error {
	// Получаем DonorResponse
	donorResponse, err := h.donorRepo.GetDonorResponseByID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor response")
	}
	if donorResponse.Status == donormodel.DonorResponseStatusCompleted {
		return apperrors.BadRequest("donor response status is invalid").WithMessage("cannot cancel a completed donation")
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.Cancel(txCtx, resID); err != nil {
			return apperrors.Internal(err, "failed to cancel donation")
		}

		// Пересчитываем статус заявки после отмены отклика
		bloodReq, err := h.bloodRepo.GetByApplicationID(txCtx, resID)
		if err != nil {
			return apperrors.Internal(err, "failed to get blood request after cancel")
		}
		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return apperrors.Internal(err, "failed to update blood request status after cancel")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
