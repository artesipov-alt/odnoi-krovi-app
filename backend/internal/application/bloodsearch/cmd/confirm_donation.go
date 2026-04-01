package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type ConfirmDonationHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	txManager *presistance.TxManager
	publisher ports.EventPublisher
}

func NewConfirmDonationHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
) *ConfirmDonationHandler {
	return &ConfirmDonationHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		txManager: txManager,
		publisher: publisher,
	}
}

func (h *ConfirmDonationHandler) Handle(ctx context.Context, donorResponseID string, factAmount int32) error {
	bloodReq, err := h.bloodRepo.GetByApplicationID(ctx, donorResponseID)
	if err != nil {
		return err
	}
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.Confirm(txCtx, donorResponseID, factAmount); err != nil {
			return err
		}

		isDone := true
		for _, resp := range bloodReq.DonorApplications {
			if !(resp.Status == model.DonorResponseStatusCompleted && resp.IsConfirmed) {
				isDone = false
				break
			}
		}

		if isDone {
			bloodReq.Close()
		}

		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return err
		}
		return nil
	})

	return err
}
