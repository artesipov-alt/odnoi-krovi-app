package query

import (
	"context"
)

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
	stats, err := h.repo.GetPortalStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
