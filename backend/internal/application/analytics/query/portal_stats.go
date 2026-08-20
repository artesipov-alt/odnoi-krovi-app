package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/middleware"
)

type PetTypeStats struct {
	Type                string
	TotalPets           int64
	ActiveBloodRequests int64
	TotalDonations      int64
	CompletedDonations  int64
	Searches            int64
	SearchVolume        float64
	DonationVolume      float64
}

type PortalStats struct {
	TotalUsers             int64
	UsersWithPhone         int64
	PhoneConversionPercent float64
	VerifiedUsers          int64
	UnverifiedUsers        int64
	TotalPets              int64
	ActiveBloodRequests    int64
	TotalDonations         int64
	CompletedDonations     int64
	TotalSearches          int64
	TotalSearchVolume      float64
	TotalDonationVolume    float64
	CatStats               PetTypeStats
	DogStats               PetTypeStats
}

// интерфейс тут же, рядом с хендлером
type PortalStatsRepository interface {
	GetPortalStats(ctx context.Context) (*PortalStats, error)
}

type PortalStatsHandler struct {
	repo PortalStatsRepository
}

func NewPortalStatsHandler(repo PortalStatsRepository) *PortalStatsHandler {
	return &PortalStatsHandler{
		repo: repo,
	}
}

func (h *PortalStatsHandler) Handle(ctx context.Context) (*PortalStats, error) {
	if middleware.GetUserRole(ctx) != "admin" {
		return nil, apperrors.Forbidden("У вас неподходящая роль для выполнения")
	}
	stats, err := h.repo.GetPortalStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
