package domain

import (
	"context"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	partnermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/partner/model"
)

type PartnerRepository interface {
	CreatePartner(ctx context.Context, name, apiKey, role, status, description string) (*partnermodel.Partner, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*partnermodel.Partner, error)
	GetByID(ctx context.Context, id string) (*partnermodel.Partner, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateLastUsedAt(ctx context.Context, id string, lastUsedAt time.Time) error
	Delete(ctx context.Context, id string) error
	ExistsByAPIKey(ctx context.Context, apiKey string) (bool, error)
	GetPartnerIdentities(ctx context.Context, partnerID string) ([]*authmodel.Identity, error)
}
