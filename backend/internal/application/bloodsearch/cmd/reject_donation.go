package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type RejectDonationHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
	txManager *presistance.TxManager
	publisher ports.EventPublisher
}

func NewRejectDonationHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
	txManager *presistance.TxManager,
	publisher ports.EventPublisher,
) *RejectDonationHandler {
	return &RejectDonationHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
		txManager: txManager,
		publisher: publisher,
	}
}

func (h *RejectDonationHandler) Handle(ctx context.Context, donorResponseID string, rejectedReason string) error {
	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}

	if err := application.Reject(rejectedReason); err != nil {
		return err
	}

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.Reject(txCtx, application); err != nil {
			return err
		}
		bloodReq, err := h.bloodRepo.GetByApplicationID(txCtx, donorResponseID)
		if err != nil {
			return err
		}
		bloodReq.RecalculateBloodAmount()
		bloodReq.RecalculateStatus()
		if err := h.bloodRepo.UpdateStatus(txCtx, bloodReq.ID, bloodReq.Status); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
