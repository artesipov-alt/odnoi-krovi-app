package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
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

func (h *RejectDonationHandler) Handle(ctx context.Context, donorResponseID string) error {
	application, err := h.donorRepo.GetDonorResponseByID(ctx, donorResponseID)
	if err != nil {
		return err
	}
	bloodReq, err := h.bloodRepo.GetByID(ctx, application.RequestID)
	if err != nil {
		return err
	}

	bloodReq.UnReserveVolume(application.Amount)

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.donorRepo.UpdateDonorResponseStatus(txCtx, donorResponseID, donormodel.DonorResponseStatusRejected); err != nil {
			return err
		}
		if err := h.bloodRepo.UpdateReservedVolume(txCtx, bloodReq.ID, bloodReq.BloodVolumeReserved); err != nil {
			return err
		}
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
