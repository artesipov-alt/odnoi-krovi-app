package job

import (
	"context"
	"log/slog"
	"time"

	bloodcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
)

type AutoConfirmJob struct {
	donorRepo              donor.Repository
	confirmDonationHandler *bloodcmd.ConfirmDonationHandler
}

func NewAutoConfirmJob(repo donor.Repository, handler *bloodcmd.ConfirmDonationHandler) *AutoConfirmJob {
	return &AutoConfirmJob{donorRepo: repo, confirmDonationHandler: handler}
}

func (j *AutoConfirmJob) Run(ctx context.Context) {
	cutoffTime := time.Now().Add(-72 * time.Hour)
	responses, err := j.donorRepo.FindNotConfirmed(ctx, cutoffTime)
	if err != nil {
		slog.Error("Failed to find not confirmed donor responses", "error", err)
		return
	}
	for _, response := range responses {
		if err := j.confirmDonationHandler.Handle(ctx, response.ID, response.Amount); err != nil {
			slog.Error("Failed to confirm donation", "error", err, "responseID", response.ID)
			continue
		}
	}
}
