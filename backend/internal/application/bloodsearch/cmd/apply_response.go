package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type ApplyResponseHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	donorRepo donor.Repository
}

func NewApplyResponseHandler(
	bloodRepo bloodsearch.BloodRequestRepository,
	donorRepo donor.Repository,
) *ApplyResponseHandler {
	return &ApplyResponseHandler{
		bloodRepo: bloodRepo,
		donorRepo: donorRepo,
	}
}

func (h *ApplyResponseHandler) Handle(ctx context.Context, donorResponseID string) error {
	req, err := h.bloodRepo.GetByApplicationID(ctx, donorResponseID)
	if err != nil {
		return err
	}
	if req.BloodVolumeReserved >= req.BloodVolumeNeeded {
		req.Close()
	}
	if err := h.donorRepo.UpdateDonorResponseStatus(ctx, donorResponseID, donormodel.DonorResponseStatusAccepted); err != nil {
		return err
	}

	return nil
}
