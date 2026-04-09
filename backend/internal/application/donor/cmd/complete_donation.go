package cmd

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type CompleteDonationHandler struct {
	donorRepo donor.Repository
}

func NewCompleteDonationHandler(
	donorRepo donor.Repository,
) *CompleteDonationHandler {
	return &CompleteDonationHandler{
		donorRepo: donorRepo,
	}
}

func (h *CompleteDonationHandler) Handle(ctx context.Context, resID string, amount float64) error {
	// Получаем DonorResponse
	donorResponse, err := h.donorRepo.GetDonorResponseByID(ctx, resID)
	if err != nil {
		return apperrors.Internal(err, "failed to get donor response")
	}
	if donorResponse.Status != donormodel.DonorResponseStatusAccepted {
		return apperrors.BadRequest("donor response status is invalid").WithMessage("donor response must be accepted to complete donation")
	}
	if err := h.donorRepo.Complete(ctx, donorResponse.ID, amount); err != nil {
		return err
	}

	return nil
}
