package query

import (
	"context"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type ListRequestsHandler struct {
	bloodRepo bloodsearch.BloodRequestRepository
	petRepo   pet.Repository
}

func NewListRequestsHandler(bloodRepo bloodsearch.BloodRequestRepository, petRepo pet.Repository) *ListRequestsHandler {
	return &ListRequestsHandler{
		bloodRepo: bloodRepo,
		petRepo:   petRepo,
	}
}

func (h *ListRequestsHandler) Handle(ctx context.Context, userID string, filters donormodel.DonorPreloadFilter) ([]*donormodel.Recipient, error) {
	pets, err := h.petRepo.GetByUserID(ctx, userID, pet.PetPreloadOptions{
		WithAll:      true,
		WithBloodReq: true,
	})
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	for i, _ := range pets {
		pets[i].RecalculateFactors(time.Now())
		pets[i].CalculateDonorStatus()
		slog.Info("", "petStatus", pets[i].PetStatus)
		slog.Info("", "petName", pets[i].StopFactors)
	}

	potentialDonors := petmodel.FilterDonors(pets)
	slog.Info("", "potentialDonors", potentialDonors)

	requests, err := h.bloodRepo.AdptiveList(ctx, potentialDonors, filters)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to list blood requests")
	}

	return requests, nil
}
