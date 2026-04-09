package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type CancelDonationHandler struct {
	donorRepo donor.Repository
}

func NewCancelDonationHandler(
	donorRepo donor.Repository,
) *CancelDonationHandler {
	return &CancelDonationHandler{
		donorRepo: donorRepo,
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

	if err := h.donorRepo.Cancel(ctx, resID); err != nil {
		return apperrors.Internal(err, "failed to cancel donation")
	}
	return nil
}
